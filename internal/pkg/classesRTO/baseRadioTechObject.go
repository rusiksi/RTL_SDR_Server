package classesRTO

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

//ЦМР для распаковки вещественных данных
const constValueLSB float32 = 100.0

//ЦМР для распаковки значения долготы
var constLonLSB float64 = 360.0 / math.Pow(2, 31)

//ЦМР для распаковки значения широты
var constLatLSB float64 = 180.0 / math.Pow(2, 31)

type BaseRTO struct {
	icao      uint32    // ICAO самолёта
	flight    string    // наименование рейса
	altitude  float32   // высота
	speed     float32   // скорость
	course    float32   // курс
	latitude  float64   // широта
	longitude float64   // долгота
	seen      int64     // время последней регистрации в эфире в UTC
	dateTime  time.Time // время последней регистрации в эфире
	messages  uint32    // количество сообщений
}

//кадр данных содержит следующую информацию
// uint32_t icao   		- ICAO
// char flight[9]   	- номер рейса, 9 байт
// uint32_t altitude    - высота
// uint32_t speed       - скорость
// uint32_t course      - курс
// int32_t lat          - широта
// int32_t lon          - долгота
// int64_t seen         - время прихода последнего пакета в UTC
// uint32_t messages    - количество сообщений
func (o *BaseRTO) Unserialize(data []byte) {
	var startByteIndex int = 0
	var sizeInt32 = 4
	var sizeInt64 = 8
	var sizeFlightStr = 9
	var offset = 0

	if len(data) <= 0{
		return
	}

	// ICAO - 0 : 4
	startByteIndex, offset = 0, sizeInt32
	o.icao = binary.LittleEndian.Uint32(data[startByteIndex:offset])

	//flight - 4:13
	startByteIndex, offset = offset, offset+sizeFlightStr
	o.flight = string(data[startByteIndex:offset])

	//altitude - 13:17
	startByteIndex, offset = offset, offset+sizeInt32
	o.altitude = float32(binary.LittleEndian.Uint32(data[startByteIndex:offset])) / constValueLSB

	//speed - 17:21
	startByteIndex, offset = offset, offset+sizeInt32
	o.speed = float32(binary.LittleEndian.Uint32(data[startByteIndex:offset])) / constValueLSB

	//course - 21:25
	startByteIndex, offset = offset, offset+sizeInt32
	o.course = float32(binary.LittleEndian.Uint32(data[startByteIndex:offset])) / constValueLSB

	//latitude - 25:29
	startByteIndex, offset = offset, offset+sizeInt32
	o.latitude = float64(binary.LittleEndian.Uint32(data[startByteIndex:offset])) * constLatLSB

	//longitude - 29 :33
	startByteIndex, offset = offset, offset+sizeInt32
	o.longitude = float64(binary.LittleEndian.Uint32(data[startByteIndex:offset])) * constLonLSB

	//seen - 33 : 41
	startByteIndex, offset = offset, offset+sizeInt64
	o.seen = int64(binary.LittleEndian.Uint64(data[startByteIndex:offset]))

	o.dateTime = time.Unix(o.seen/1000, 0)

	//messages - 41 : 45
	startByteIndex, offset = offset, offset+sizeInt32
	o.messages = binary.LittleEndian.Uint32(data[startByteIndex:offset])
}

func (o *BaseRTO) String() string {
	return fmt.Sprintf("aircraft %X (flyight %s) \n"+
		"\t altitude %f\n"+
		"\t speed %f\n"+
		"\t course %f\n"+
		"\t latitude %f\n"+
		"\t longitude %f\n"+
		"\t seen %d\n"+
		"\t time %v\n"+
		"\t message %d\n",
		o.icao,
		o.flight,
		o.altitude,
		o.speed,
		o.course,
		o.latitude,
		o.longitude,
		o.seen,
		o.dateTime,
		o.messages)
}
