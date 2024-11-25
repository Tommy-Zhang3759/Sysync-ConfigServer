package APIWorkers

import (
	"ConfigServer/clientManage"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestSearchNewClient_Run(t *testing.T) {
	clientManage.Init("../../data/clientInfo.db")
	p := NewConnNewClient(clientManage.Container)

	p.SetKeyWord("key")
	_ = clientManage.CliUdpApiGateway.Add(p)

	go func() {
		_ = clientManage.CliUdpApiGateway.Run()
	}()

	_ = p.Start()

	type Message struct {
		FName     string `json:"f_name"`
		Mac       string `json:"mac"`
		OSVersion string `json:"os_version"`
		ProductID string `json:"product_id"`
		HostName  string `json:"host_name"`
		Status    int    `json:"status_code"`
	}

	var message []byte

	message, _ = json.Marshal(Message{
		FName:     "key",
		Mac:       "11:11:11:11:11:11",
		OSVersion: "12345567",
		ProductID: "09876",
		HostName:  "TEST-HOST",
		Status:    0,
	})

	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6004,
	}

	for range 1000 {
		go func() {
			conn, _ := net.DialUDP("udp", nil, addr)
			defer func(conn *net.UDPConn) {
				_ = conn.Close()
			}(conn)

			_, _ = conn.Write(message)
		}()
	}

	time.Sleep(10 * time.Second)

}
