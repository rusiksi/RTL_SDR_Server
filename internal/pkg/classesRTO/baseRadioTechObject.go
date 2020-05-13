package classesRTO

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

const constValueLSB float32 = 100.0

var constLonLSB float64 = 360.0 / math.Pow(2, 31)

var constLatLSB float64 = 180.0 / math.Pow(2, 31)

type BaseRTO struct {
	icao      uint32
	flight    string
	altitude  float32
	speed     float32
	course    float32
	latitude  float64
	longitude float64
	seen      uint64
	dateTime  time.Time
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
func (o *BaseRTO) Unserialize(data []byte) {
	o.icao = binary.LittleEndian.Uint32(data[0:4])
	o.flight = string(data[4:13])
	o.altitude = float32(binary.LittleEndian.Uint32(data[13:17])) / constValueLSB
	o.speed = float32(binary.LittleEndian.Uint32(data[17:21])) / constValueLSB
	o.course = float32(binary.LittleEndian.Uint32(data[21:25])) / constValueLSB
	o.latitude = float64(binary.LittleEndian.Uint32(data[25:29])) * constLatLSB
	o.longitude = float64(binary.LittleEndian.Uint32(data[29:33])) * constLonLSB
	o.dateTime = time.Unix(int64(binary.LittleEndian.Uint64(data[33:41])),0)
	o.seen = binary.LittleEndian.Uint64(data[33:41])
	o.messages = binary.LittleEndian.Uint32(data[41:45])
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