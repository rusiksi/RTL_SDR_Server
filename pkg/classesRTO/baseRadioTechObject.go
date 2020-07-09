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
	Icao      uint32    `json:"icao"`      // ICAO самолёта
	Flight    string    `json:"flight"`    // наименование рейса
	Altitude  float32   `json:"altitude"`  // высота
	Speed     float32   `json:"speed"`     // скорость
	Course    float32   `json:"course"`    // курс
	Latitude  float64   `json:"latitude"`  // широта
	Longitude float64   `json:"longitude"` // долгота
	Seen      int64     `json:"seen"`      // время последней регистрации в эфире в UTC
	DateTime  time.Time `json:"dateTime"`  // время последней регистрации в эфире
	Messages  uint32    `json:"messages"`  // количество сообщений
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

	if len(data) <= 0 {
		return
	}

	// ICAO - 0 : 4
	startByteIndex, offset = 0, sizeInt32
	o.Icao = binary.LittleEndian.Uint32(data[startByteIndex:offset])

	//flight - 4:13
	startByteIndex, offset = offset, offset+sizeFlightStr
	o.Flight = string(data[startByteIndex:offset])

	//altitude - 13:17
	startByteIndex, offset = offset, offset+sizeInt32
	o.Altitude = float32(binary.LittleEndian.Uint32(data[startByteIndex:offset])) / constValueLSB

	//speed - 17:21
	startByteIndex, offset = offset, offset+sizeInt32
	o.Speed = float32(binary.LittleEndian.Uint32(data[startByteIndex:offset])) / constValueLSB

	//course - 21:25
	startByteIndex, offset = offset, offset+sizeInt32
	o.Course = float32(binary.LittleEndian.Uint32(data[startByteIndex:offset])) / constValueLSB

	//latitude - 25:29
	startByteIndex, offset = offset, offset+sizeInt32
	o.Latitude = float64(binary.LittleEndian.Uint32(data[startByteIndex:offset])) * constLatLSB

	//longitude - 29 :33
	startByteIndex, offset = offset, offset+sizeInt32
	o.Longitude = float64(binary.LittleEndian.Uint32(data[startByteIndex:offset])) * constLonLSB

	//seen - 33 : 41
	startByteIndex, offset = offset, offset+sizeInt64
	o.Seen = int64(binary.LittleEndian.Uint64(data[startByteIndex:offset]))

	o.DateTime = time.Unix(o.Seen/1000, 0)

	//messages - 41 : 45
	startByteIndex, offset = offset, offset+sizeInt32
	o.Messages = binary.LittleEndian.Uint32(data[startByteIndex:offset])
}

func (o *BaseRTO) Serialize() []byte {
	var startByteIndex int = 0
	var sizeInt32 = 4
	var sizeInt64 = 8
	var sizeFlightStr = 9
	var offset = 0

	data := make([]byte, 46)

	// ICAO - 0 : 4
	startByteIndex, offset = 0, sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], o.Icao)

	//flight - 4:13
	startByteIndex, offset = offset, offset+sizeFlightStr
	b := make([]byte, 9)
	copy(b, o.Flight)

	copy(data[startByteIndex:offset], b)

	//altitude - 13:17
	startByteIndex, offset = offset, offset+sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], uint32(o.Altitude*constValueLSB))

	//speed - 17:21
	startByteIndex, offset = offset, offset+sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], uint32(o.Speed*constValueLSB))

	//course - 21:25
	startByteIndex, offset = offset, offset+sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], uint32(o.Course*constValueLSB))

	//latitude - 25:29
	startByteIndex, offset = offset, offset+sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], uint32(o.Latitude/constLatLSB))

	//longitude - 29 :33
	startByteIndex, offset = offset, offset+sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], uint32(o.Longitude/constLonLSB))

	//seen - 33 : 41
	startByteIndex, offset = offset, offset+sizeInt64
	binary.LittleEndian.PutUint64(data[startByteIndex:offset], uint64(o.Seen))

	//messages - 41 : 45
	startByteIndex, offset = offset, offset+sizeInt32
	binary.LittleEndian.PutUint32(data[startByteIndex:offset], o.Messages)

	return data
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
		o.Icao,
		o.Flight,
		o.Altitude,
		o.Speed,
		o.Course,
		o.Latitude,
		o.Longitude,
		o.Seen,
		o.DateTime,
		o.Messages)
}
