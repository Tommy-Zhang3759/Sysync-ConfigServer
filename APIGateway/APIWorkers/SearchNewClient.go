package APIWorkers

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"

	"encoding/json"
	"fmt"
	"log"
	"net"
)

type ConnNewClient struct {
	APIGateway.UDPAPIPortTemp
	cliContainer *clientManage.CliContainer
}

func (u *ConnNewClient) Run() error {
	stop := false
	// TODO: use api for data reading
	//var db DataFrame.DataFrame = &DataFrame.SQLite{}
	//
	//err := db.Connect("data/clientInfo.db")
	//if err != nil {
	//	log.Fatalf("Failed to connect to database: %v", err)
	//}
	//defer func() {
	//	if err = db.Close(); err != nil {
	//		log.Fatalf("Failed to disconnect to database: %v", err)
	//	}
	//}()

	for stop == false {
		reqPack := u.MessageQue.Pop().(APIGateway.UDPMessage)

		select {
		case <-u.EndRun:
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:

			type ClientRequest struct {
				Mac       string `json:"mac"`
				OSVersion string `json:"os_version"`
				ProductID string `json:"product_id"`
				HostName  string `json:"host_name"`
			}

			type rspTemp struct {
				APIGateway.ApiResponse
				Host     string `json:"host_name,omitempty"`
				IP       string `json:"host_ip,omitempty"`
				SysyncID string `json:"sysync_id,omitempty"`
			}

			var rsp rspTemp
			rsp.Fname = u.GetKeyWord()

			var hostname string = reqPack.Text["host_name"].(string)
			var ip net.IP = reqPack.Addr.IP
			var mac, err = net.ParseMAC(reqPack.Text["mac"].(string))
			var osVersion string = reqPack.Text["os_version"].(string)
			var productID string = reqPack.Text["product_id"].(string)
			var status = int(reqPack.Text["status_code"].(float64))
			if err != nil {
				rsp.Error = err.Error()
				rsp.Status = -1
			} else {
				newCli := clientManage.CreateNewClientInfo(
					hostname,
					ip,
					mac,
					status,
					osVersion,
					productID,
				)
				err = clientManage.Push(newCli)
				if err != nil {
					log.Println(err)
				}
			}

			rspJson, err := json.Marshal(rsp)
			if err != nil {
				log.Printf("Error marshalling response: %v", err)
			}

			err = u.Gateway.SendMess(rspJson, reqPack.Addr)
			if err != nil {
				log.Printf("Error sending request: %v", err)
			}
		}
	}
	return nil
}
