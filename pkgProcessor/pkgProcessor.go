package pkgProcessor

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)


const constValueLSB float32 = 100.0
var constLonLSB float64 = 360.0 / math.Pow(2, 31);
var constLatLSB float64 = 180.0 / math.Pow(2, 31);

type objectInfo struct {
	icao      uint32
	flight    string
	altitude  float32
	speed     float32
	course    float32
	latitude  float64
	longitude float64
	seen      uint64
	messages  uint32
}

/* кадр данных содержит следующую информацию
 * <br> uint32_t icao   - ICAO</br>
 * <br> char flight[9]   - номер рейса, 9 байт</br>
 * <br> uint32_t altitude    - высота</br>
 * <br> uint32_t speed       - скорость</br>
 * <br> uint32_t course      - курс</br>
 * <br> int32_t lat          - широта</br>
 * <br> int32_t lon          - долгота</br>
 * <br> int64_t seen         - время прихода последнего пакета в UTC</br>
 * <br> uint32_t messages    - количество сообщений</br>
 */
func (o* objectInfo) unserialize(data []byte)  {
	o.icao = binary.LittleEndian.Uint32(data[0 : 4])
	o.flight = string(data[4 : 13])
	o.altitude = float32(binary.LittleEndian.Uint32(data[13 : 17])) / constValueLSB
	o.speed  = float32(binary.LittleEndian.Uint32(data[17 : 21])) / constValueLSB
	o.course = float32(binary.LittleEndian.Uint32(data[21 : 25])) / constValueLSB
	o.latitude = float64(binary.LittleEndian.Uint32(data[25 : 29])) * constLatLSB
	o.longitude = float64(binary.LittleEndian.Uint32(data[29 : 33])) * constLonLSB
	o.seen = binary.LittleEndian.Uint64(data[33 : 41])
	o.messages = binary.LittleEndian.Uint32(data[41 : 45])
}

func (o* objectInfo) String() string {
	return fmt.Sprintf("aircraft %X (flyight %s) \n" +
		"\t altitude %f\n" +
	"\t speed %f\n" +
		"\t course %f\n" +
		"\t latitude %f\n" +
		"\t longitude %f\n" +
		"\t seen %d\n" +
		"\t message %d\n",
		o.icao,
		o.flight,
		o.altitude,
		o.speed,
		o.course,
		o.latitude,
		o.longitude,
		o.seen,
		o.messages)
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

func ParseRawData(ch <-chan []byte) {
	var sizeFrame uint32 = 0
	var countFrame uint32 = 0
	var object = new(objectInfo)

	for sl := range ch {
		sizeFrame = binary.LittleEndian.Uint32(sl[:4])
		countFrame = binary.LittleEndian.Uint32(sl[4:8])

		fmt.Println("sizeFrame = ", sizeFrame, "countFrame = ", countFrame)
		var i uint32
		for i = 0; i < countFrame; i++ {
			object.unserialize(sl[8 + sizeFrame * i :])
			fmt.Println(object)
		}

		printerRawData(sl)

	}
}

func printerRawData(data []byte) {
	now := time.Now()
	fmt.Printf("==== %v ====\n", now)
	fmt.Printf("%s", hex.Dump(data))
	fmt.Println("==========================================================")
}

