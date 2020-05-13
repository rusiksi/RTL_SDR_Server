package pkgProcessor

import (
	"RTL_SDR_Server/internal/pkg/classesRTO"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)
type IPkgProcessor interface {
	ReadData(conn net.Conn) error
	processData(chan dataFrame)
}


type PkgProcessorImpl struct {
	sliceFrame []dataFrame

}

type dataFrame struct {
	sizeFrame  uint32
	countFrame uint32
	data       []byte
}

const constMinSize int = 8

func (pkgProc* PkgProcessorImpl) ReadData(conn net.Conn) error  {

	if pkgProc == nil {
		return errors.New(" nil pointer in ReadData\n")
	}

	msgHdr := dataFrame{}
	//чтение заголовка пакета
	buffer := make([]byte, constMinSize)
	hdr := bytes.NewReader(buffer)

	ch := make(chan dataFrame)

	go pkgProc.processData(ch)

	for {
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

		hdr.Seek(0,0)
		err = binary.Read(hdr, binary.LittleEndian, &msgHdr.sizeFrame)
		err = binary.Read(hdr, binary.LittleEndian, &msgHdr.countFrame)

		if err != nil  || msgHdr.countFrame == 0 ||  msgHdr.sizeFrame == 0{
			return errors.New("Error reading data size \n")
		}

		var sizePkg = int(msgHdr.countFrame * msgHdr.sizeFrame)

		if len(msgHdr.data) < sizePkg{
			msgHdr.data = make([]byte, sizePkg)
		}

		//подготовка массива для чтения
		lenData, err := conn.Read(msgHdr.data)

		if lenData != sizePkg {
			return errors.New("Error reading main data block \n")
		}


		ch <- msgHdr
	}
	close(ch)
	return nil
}

func (pkgProc* PkgProcessorImpl) processData( df chan dataFrame){
	if pkgProc == nil {
		return
	}
	for v:= range df {
		pkgProc.parseRawData(&v)
	}
}

/**
* Распаковка пришедшего массива данных.
* в массиве данные упакованы следующим образом :
* <br>==========Преамбула========</br>
* <br>uint32_t sizeFrame   - размер 1 блока данных</br>
* <br>uint32_t countFrame  - количество блоков данных</br>
*
* <br>==========1...countFrame кадров данных========</br>
 *
**/
func (pkgProc* PkgProcessorImpl) parseRawData(df* dataFrame) {

		fmt.Println("sizeFrame = ", df.sizeFrame, "countFrame = ", df.countFrame)
		var i uint32
		for i = 0; i < df.countFrame; i++ {
			var object = new(classesRTO.BaseRTO)
			object.Unserialize(df.data[df.sizeFrame*i:])
			fmt.Println(object)
		}
}
