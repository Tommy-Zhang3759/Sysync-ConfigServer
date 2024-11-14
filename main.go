package main

import (
	"ConfigServer/clientManage"
	"ConfigServer/webUI"
	"time"
)

var Version string

func init() {

}

func main() {
	//var mvg sync.WaitGroup

	clientManage.Init("data/clientInfo.db")

	go func() {
		_ = clientManage.CliUdpApiGateway.Run()
	}()

	go webUI.StartServer("8080", webUI.Handler)

	//mvg.Wait()
	// var clients []clientManage.Client
	for {
		time.Sleep(time.Hour)
	}
}
