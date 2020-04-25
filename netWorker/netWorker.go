package netWorker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type msgHeader struct{
	SizeFrame uint32
	CountFrame uint32
}

const constMinSize int = 8


type bufferedConn struct {
	r        *bufio.Reader
	net.Conn
}


func newBufferedConn(c net.Conn) bufferedConn {
	return bufferedConn{bufio.NewReader(c), c}
}

func newBufferedConnSize(c net.Conn, n int) bufferedConn {
	return bufferedConn{bufio.NewReaderSize(c, n), c}
}

func (b bufferedConn) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b bufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

func (b bufferedConn) ReadBytes(p []byte) (int, error) {
	return b.r.Read(p)
}

func HandleConnection(conn net.Conn, ch chan  <- [] byte) {
	name := conn.RemoteAddr().String()

	fmt.Printf("%+v connected\n", name)


	defer conn.Close()
	reader := newBufferedConnSize(conn,constMinSize)

	var msgHdr msgHeader
	for {
		//полчение заголовка
		b, err := reader.Peek(constMinSize)
		if err != nil {
			break
		}
		hdr := bytes.NewReader(b)
		err = binary.Read(hdr, binary.LittleEndian, &msgHdr)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
			break
		}
		var sizePkg = int(msgHdr.CountFrame * msgHdr.SizeFrame) + constMinSize
		//подготовка массива для чтения
		bufData := make([]byte,sizePkg)

		var lenRead = 0
		for lenRead < sizePkg{
			t := make([]byte,sizePkg,sizePkg)
			len, err := conn.Read(t)

			if len < 0 || err != nil{
				fmt.Println("binary.Read failed:", err)
				break
			}
			copy(bufData[lenRead:],t[:len])
			lenRead += len
		}
		if lenRead != sizePkg {
			fmt.Println("Lost part of data")
		}
		ch <- bufData
	}
	close(ch)
}
