package APIWorkers

import (
	"ConfigServer/APIGateway"
	"encoding/json"
	"net"
	"testing"
	"time"
)

func TestSearchNewClient_Run_Run(t *testing.T) {
	var gateway = APIGateway.NewUDPAPIGateway(6004, "0.0.0.0")
	p := SearchNewClient{}
	p.KeyWord = "key"

	_ = gateway.Init()

	go func() {
		_ = gateway.Run()
	}()
	_ = gateway.Add(&p)
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
	}

	var message []byte

	message, _ = json.Marshal(Message{
		FName:     "key",
		Mac:       "11:11:11:11:11",
		OSVersion: "12345567",
		ProductID: "09876",
		HostName:  "TEST-HOST",
	})

	// 定义本地的目标地址（本地 IP 地址 127.0.0.1 和端口 6004）
	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6004,
	}

	conn, _ := net.DialUDP("udp", nil, addr)
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	// 发送数据
	_, _ = conn.Write(message)

}
