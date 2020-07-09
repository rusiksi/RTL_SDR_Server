package pkgProcessor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/art-injener/RTL_SDR_Server/configs"
	"github.com/art-injener/RTL_SDR_Server/pkg/classesRTO"
	"github.com/streadway/amqp"
	"io"
	"log"
	"net"
)

const CH_SIZE int = 100
const constMinSize int = 8 // минимальный размер пакета, котором передаётся размер и количество байт

// интерфейс для реализации процессора данных
// ReadData - чтение из сетевого интерфеса блока данных
// processData - функция для обработки бинарных пакетов данных,поступающих через канал
type IPkgProcessor interface {
	ReadData(conn net.Conn) error
	processData(chan dataFrame)
}

// структура для хранения и передачи блока данных
type dataFrame struct {
	sizeFrame  uint32 // размер блока данных
	countFrame uint32 // количество блоков данных в пакете
	data       []byte // бинарный массиов данных
}

// имплеоментация интерфеса процессора IPkgProcessor
type PkgProcessorImpl struct {
	sliceFrame []dataFrame
}

func (pkgProc *PkgProcessorImpl) ReadData(conn net.Conn) error {

	//чтение заголовка пакета
	buffer := make([]byte, constMinSize)
	hdr := bytes.NewReader(buffer)

	if pkgProc == nil || hdr == nil {
		return errors.New(" nil pointer in ReadData\n")
	}

	//структура для описания пакета с данными
	msgHdr := dataFrame{}
	ch := make(chan dataFrame, CH_SIZE)

	// обработка приходящих пакетов выполниется в отдельной горутине
	go pkgProc.processData(ch)

	for {
		//чтение заголовка
		lenRead, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if lenRead != constMinSize {
			return errors.New("Error reading data size \n")
		}

		hdr.Seek(0, 0)
		// получаем размер 1 фрейма данных
		err = binary.Read(hdr, binary.LittleEndian, &msgHdr.sizeFrame)
		//получаем количество фреймов данных
		err = binary.Read(hdr, binary.LittleEndian, &msgHdr.countFrame)

		if err != nil || msgHdr.countFrame == 0 || msgHdr.sizeFrame == 0 {
			return errors.New("Error reading data size \n")
		}
		//вычисление размера для чтени яданных
		var sizePkg = int(msgHdr.countFrame * msgHdr.sizeFrame)
		//увеличение буфера для чтения, если данные не помещаются
		if len(msgHdr.data) < sizePkg {
			msgHdr.data = make([]byte, sizePkg)
		}
		//подготовка массива для чтения
		lenData, err := conn.Read(msgHdr.data)

		if lenData != sizePkg {
			log.Printf("Read bytes = %d, but sizePkg = %d  \n", lenData, sizePkg)
			return errors.New("Error reading main data block \n")
		}
		// передача в канал для последующего декодирования данных
		ch <- msgHdr
	}
	close(ch)
	return nil
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}

}

// выполняется в отдельной горутине,обработка закончится,
// когда канал с данными будет закрыт
func (pkgProc *PkgProcessorImpl) processData(df chan dataFrame) {
	if pkgProc == nil {
		return
	}

	//TODO:вынести подключение в отдельную функцию
	configRMQ := configs.GetRMQConfig()
	conn, err := amqp.Dial(configRMQ.AMQPConnectionURL)

	handleError(err, "Can't connect to AMQP")

	defer conn.Close()

	amqpChannel, err := conn.Channel()

	handleError(err, "Can't create a amqpChannel")

	queue, err := amqpChannel.QueueDeclare("add", true, false, false, false, nil)
	handleError(err, "Could not declare `add` queue")
	defer amqpChannel.Close()

	for v := range df {
		pkgProc.parseRawData(v, amqpChannel, &queue)
	}
}

// Распаковка пришедшего массива данных.
// в массиве данные упакованы следующим образом :
// ==========Преамбула========
// uint32_t sizeFrame   - размер 1 блока данных
// uint32_t countFrame  - количество блоков данных
// ==========1...countFrame кадров данных========
func (pkgProc *PkgProcessorImpl) parseRawData(df dataFrame, amqpChannel *amqp.Channel, queue *amqp.Queue) {

	//fmt.Println("sizeFrame = ", df.sizeFrame, "countFrame = ", df.countFrame)
	//fmt.Println(hex.Dump(df.data))

	var i uint32
	for i = 0; i < df.countFrame; i++ {
		if len(df.data) < int(df.sizeFrame*i) {
			continue
		}
		var object = new(classesRTO.BaseRTO)

		object.Unserialize(df.data[df.sizeFrame*i:])

		if amqpChannel != nil && queue != nil {
			body, err := json.Marshal(object)
			if err != nil {
				handleError(err, "Error encoding JSON")
			}

			err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})

			if err != nil {
				handleError(err, "Error publishing message")
			}
		}
		fmt.Println(object)
	}
}