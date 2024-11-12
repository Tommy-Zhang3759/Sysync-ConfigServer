package APIWorkers

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"
	DataFrame "ConfigServer/utils/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type HostNameReq struct {
	APIGateway.UDPAPIPortTemp
}

func (u *HostNameReq) Run() error {
	stop := false

	var db DataFrame.DataFrame = clientManage.DataFrameConn()
	defer func() {
		if err := db.Close(); err != nil {
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
				Mac string `json:"mac"`
			}

			type rspTemp struct {
				FName string `json:"f_name"`
				Host  string `json:"host_name"`
				IP    string `json:"host_ip"`
			}

			rsp, err := func(macAddress string) (rspTemp, error) {
				query := "SELECT host_name, IP_address FROM win_cli WHERE MAC_address = ?"
				rows, err := db.Query(query, macAddress)

				if err != nil {
					return rspTemp{}, err
				}
				defer func(rows *sql.Rows) {
					_ = rows.Close()
				}(rows)

				if rows.Next() {
					var hostName, ipAddress string
					if err := rows.Scan(&hostName, &ipAddress); err != nil {
						return rspTemp{}, err
					}
					if rows.Next() {
						return rspTemp{}, fmt.Errorf("HostName and IP already exists")
					}

					return rspTemp{
						FName: "host_name_offer",
						Host:  hostName,
						IP:    ipAddress,
					}, nil
				}
				return rspTemp{}, sql.ErrNoRows
			}(reqPack.Text["mac"].(string))

			if err != nil {
				log.Printf("Error finding client information: %v", err)
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
