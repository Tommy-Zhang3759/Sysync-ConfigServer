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

func NewConnNewClient(cliContainer *clientManage.CliContainer) *ConnNewClient {
	return &ConnNewClient{cliContainer: cliContainer}
}

func (u *ConnNewClient) Start() error {
	go func() {
		_ = u.run()
	}()
	return nil
}

func (u *ConnNewClient) run() error {
	stop := false

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

			go func() {
				var rsp rspTemp
				rsp.Fname = u.GetKeyWord()

				var hostname string = reqPack.Text["host_name"].(string)
				var ip net.IP = reqPack.Addr.IP
				var port = reqPack.Addr.Port
				var mac, macErr = net.ParseMAC(reqPack.Text["mac"].(string))
				var osVersion string = reqPack.Text["os_version"].(string)
				var productID string = reqPack.Text["product_id"].(string)
				var status = int(reqPack.Text["status_code"].(float64))
				// TODO: enable force create to jump mac existence check
				if macErr != nil {
					rsp.Error = macErr.Error()
					rsp.Status = -1
				} else {
					newCli := clientManage.CreateNewClientInfo(
						hostname,
						ip,
						port,
						mac,
						status,
						osVersion,
						productID,
					)
					if ex, _ := clientManage.MacExists(mac.String()); !ex {
						macErr = clientManage.Push(newCli)
					} else {
						log.Printf("trying to log existed MAC: %s", mac.String())
					}

					if macErr != nil {
						log.Println(macErr)
					}
				}

				rspJson, macErr := json.Marshal(rsp)
				if macErr != nil {
					log.Printf("Error marshalling response: %v", macErr)
				}

				macErr = u.Gateway.SendMess(rspJson, reqPack.Addr)
				if macErr != nil {
					log.Printf("Error sending request: %v", macErr)
				}
			}()
		}
	}
	return nil
}
