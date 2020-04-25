package main

import (
	"RTL_SDR_Server/netWorker"
	"RTL_SDR_Server/pkgProcessor"
	"net"
)


func main() {
	listner, err := net.Listen("tcp", "192.168.0.103:62000")
	if err != nil {
		panic(err)
	}

	pkgChn := make(chan []byte, 1)
	for {
		conn, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		go netWorker.HandleConnection(conn,pkgChn)
		go pkgProcessor.ParseRawData(pkgChn)
	}
}