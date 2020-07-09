package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/art-injener/RTL_SDR_Server/pkg/classesRTO"
	"log"
	"net"
	"time"
)

const constMinSize int = 8 // минимальный размер пакета, котором передаётся размер и количество байт

type IPkgBuilder interface {
	WriteData(conn net.Conn) error
	InitImitObject()
}

// структура для хранения и передачи блока данных
type dataFrame struct {
	sizeFrame  uint32 // размер блока данных
	countFrame uint32 // количество блоков данных в пакете
	data       []byte // бинарный массиов данных
}

type PkgBuilderImpl struct {
	mapObject map[uint32]*classesRTO.BaseRTO
}

func (pkgBuilder *PkgBuilderImpl) InitImitObject() {

	pkgBuilder.mapObject = make(map[uint32]*classesRTO.BaseRTO)

	object := new(classesRTO.BaseRTO)
	object.Icao = 0xABADBABE
	object.Flight = "Imit_1"
	object.Altitude = 10000
	object.Speed = 567
	object.Course = 234
	object.Latitude = 47.231
	object.Longitude = 39.723
	object.Messages = 1
	object.Seen = time.Now().UnixNano() / int64(time.Millisecond)
	object.DateTime = time.Unix(object.Seen/1000, 0)

	pkgBuilder.mapObject[object.Icao] = object

	object = new(classesRTO.BaseRTO)
	object.Icao = 0xCAFEBABE
	object.Flight = "Imit_2"
	object.Altitude = 4598
	object.Speed = 188
	object.Course = 90
	object.Latitude = 47.251
	object.Longitude = 39.623
	object.Messages = 1
	object.Seen = time.Now().UnixNano() / int64(time.Millisecond)
	object.DateTime = time.Unix(object.Seen/1000, 0)

	pkgBuilder.mapObject[object.Icao] = object

	for k, v := range pkgBuilder.mapObject {
		fmt.Printf("%d -> %s\n", k, v.String())
	}

}

func (pkgBuilder *PkgBuilderImpl) updateImitObject() {
	for _, v := range pkgBuilder.mapObject {
		v.Messages++
		v.Seen = time.Now().UnixNano() / int64(time.Millisecond)
		v.DateTime = time.Unix(v.Seen/1000, 0)
		v.Longitude += 0.001
		v.Latitude += 0.001
	}
}

func (pkgBuilder *PkgBuilderImpl) WriteData(conn net.Conn) error {

	//структура для описания пакета с данными
	msgHdr := dataFrame{}

	var cnt uint32 = 0
	var frameSize uint32 = 0

	for _, v := range pkgBuilder.mapObject {
		cnt++
		data := v.Serialize()
		frameSize = uint32(len(data))
		msgHdr.data = append(msgHdr.data, data...)
	}

	msgHdr.countFrame = cnt
	msgHdr.sizeFrame = frameSize

	//чтение заголовка пакета
	buffer := make([]byte, constMinSize+int(frameSize)*int(cnt))

	hdr := bytes.NewBuffer(buffer)

	if pkgBuilder == nil || hdr == nil {
		return errors.New(" nil pointer in WriteData\n")
	}

	hdr.Reset()

	// получаем размер 1 фрейма данных
	err := binary.Write(hdr, binary.LittleEndian, &msgHdr.sizeFrame)
	handleError(err, "Error set size of data frame")

	//получаем количество фреймов данных
	err = binary.Write(hdr, binary.LittleEndian, &msgHdr.countFrame)
	handleError(err, "Error set number of data frame")

	err = binary.Write(hdr, binary.LittleEndian, &msgHdr.data)
	handleError(err, "Error set object's data")

	_, err = conn.Write(hdr.Bytes())

	handleError(err, "Error write data to net.")

	pkgBuilder.updateImitObject()
	return nil
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}

}
