package APIGateway

import (
	"net"
	"testing"
)

var gateway = NewUDPAPIGateway(6004)

func TestHostNameReq_Run(t *testing.T) {
	p := HostNameReq{}
	p.keyWord = "key"

	gateway.Init()
	go gateway.Run()
	_ = gateway.Add(&p)
	go func() {
		_ = p.Run()
	}()

	message := []byte("{\"mac\": \"8c:ec:4b:a4:04:2e\", \"f_name\": \"key\"}")

	// 定义本地的目标地址（本地 IP 地址 127.0.0.1 和端口 6004）
	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 6004,
	}

	// 创建一个 UDP 连接（使用 nil 作为源地址，表示没有特定的源地址）
	conn, _ := net.DialUDP("udp", nil, addr)
	defer conn.Close()

	// 发送数据
	_, _ = conn.Write(message)

}
