package main

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"
	"ConfigServer/webUI"
	"time"
)

var Version string

var UdpHostPort = 6004
var UdpClientPort = 6003

func Init(dbPath string) {
	// TODO: support identify the server ip that is under the same net range as clients
	APIGateway.CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpHostPort, "0.0.0.0")
	APIGateway.CliUdpApiGateway.Init()
	go func() {
		_ = APIGateway.CliUdpApiGateway.Run()
	}()

	clientManage.Init(dbPath)
}

func main() {
	//var mvg sync.WaitGroup

	Init("data/clientInfo.db")

	go webUI.StartServer("8080", webUI.Handler)

	//mvg.Wait()
	// var clients []clientManage.Client
	for {
		time.Sleep(time.Hour)
	}
}
