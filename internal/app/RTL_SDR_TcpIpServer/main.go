package main

import (
	"RTL_SDR_Server/internal/pkg/netWorker"
	"log"
)

func main() {

 	config := netWorker.NewConfig()
	tcpServer, err := netWorker.NewServer(config)
	if err != nil {
		panic(err)
	}

	err = tcpServer.Run()
	if err != nil {
		log.Fatalf("Server stops with error : %v ",err.Error())
	}
}
