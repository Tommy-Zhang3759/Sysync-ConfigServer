package APIWorkers

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"
	DataFrame "ConfigServer/utils/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type SearchNewClient struct {
	APIGateway.UDPAPIPortTemp
}

func (u *SearchNewClient) Run() error {
	stop := false

	var db DataFrame.DataFrame = &DataFrame.SQLite{}

	err := db.Connect("data/clientInfo.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("Failed to disconnect to database: %v", err)
		}
	}()

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
			rsp.Fname = u.KeyWord

			var hostname string = reqPack.Text["host_name"].(string)
			var ip net.UDPAddr = reqPack.Addr
			var mac, err = net.ParseMAC(reqPack.Text["mac"].(string))
			var osVersion string = reqPack.Text["os_version"].(string)
			var productID string = reqPack.Text["product_id"].(string)
			var status = reqPack.Text["status_code"].(int)
			if err != nil {
				rsp.Error = err.Error()
				rsp.Status = -1
			} else {
				_ = clientManage.CreateNewClientInfo(
					hostname,
					&ip,
					mac,
					status,
					osVersion,
					productID,
				)
				func() {
					query := "SELECT host_name, IP_address, sysync_ID FROM win_cli WHERE MAC_address = ? "
					rows, err := db.Query(query, mac)
					defer func(rows *sql.Rows) {
						_ = rows.Close()
					}(rows)
					if err != nil {
						log.Fatalf("Failed to query win_cli: %v", err)
					}

					if rows.Next() {
						var hostName, ipAddress string
						var sysyncID []byte
						_ = rows.Scan(&hostName, &ipAddress, sysyncID)
						rsp.Status = -1
						rsp.Error = "MAC already exists, logged under " + hostName + ipAddress
					}
				}()
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
