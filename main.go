package main

import (
	"ConfigServer/webUI"
	"time"
)

func init() {

}

func main() {
	//var mvg sync.WaitGroup

	go webUI.StartServer("8080", webUI.Handler)

	//mvg.Wait()
	// var clients []clientManage.Client
	for {
		time.Sleep(time.Hour)
	}
}
