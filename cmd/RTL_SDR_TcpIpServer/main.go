package main

import (
	"github.com/art-injener/RTL_SDR_Server/configs"
	"github.com/art-injener/RTL_SDR_Server/pkg/netWorker"
	"log"
)

func main() {

	config := configs.GetNetworkConfig()
	tcpServer, err := netWorker.NewServer(config)
	if err != nil {
		panic(err)
	}

	err = tcpServer.Run()
	if err != nil {
		log.Fatalf("Server stops with error : %v ", err.Error())
	}
}
