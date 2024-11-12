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
	p := ConnNewClient{
		cliContainer: clientManage.Container,
	}
	p.SetKeyWord("key")
	clientManage.CliUdpApiGateway.Add(&p)

	go func() {
		_ = clientManage.CliUdpApiGateway.Run()
	}()

	go func() {
		_ = p.Run()
	}()

	time.Sleep(2 * time.Second)

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

	conn, _ := net.DialUDP("udp", nil, addr)
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	_, _ = conn.Write(message)

	for true {
		time.Sleep(time.Second)
	}

}