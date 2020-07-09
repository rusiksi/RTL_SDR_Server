package main

import (
	"github.com/art-injener/RTL_SDR_Server/configs"
	"log"
	"net"
	"time"
)

func main() {

	config := configs.GetNetworkConfig()

	con, err := net.Dial(config.Protocol, "server:62000")

	if err != nil {
		log.Fatal("could not connect to server: ", err)
		return
	}

	log.Printf("Connect to %s:%s",
		con.RemoteAddr().Network(),
		con.RemoteAddr().String())

	defer con.Close()

	var pkgBuilder IPkgBuilder = new(PkgBuilderImpl)
	pkgBuilder.InitImitObject()
	for {
		pkgBuilder.WriteData(con)
		time.Sleep(1 * time.Second)
	}
}
