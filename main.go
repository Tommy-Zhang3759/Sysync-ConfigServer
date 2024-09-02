package main

import (
	"ConfigServer/clientManage"
	"ConfigServer/webUI"
	"time"
)

func init() {

}

func main() {
	//var mvg sync.WaitGroup

	go webUI.StartServer("8080", webUI.Handler)

	go clientManage.APIService("6004")

	//mvg.Wait()
	// var clients []clientManage.Client
	for {
		time.Sleep(time.Hour)
	}
}
