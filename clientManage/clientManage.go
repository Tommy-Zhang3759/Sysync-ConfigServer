package clientManage

import (
	"fmt"
	"math/rand"
	"net"
)

type Client struct {
	addr   net.Addr
	conn   *net.Conn
	mac    net.HardwareAddr
	status int // stat code
	id     string
}

func NewClient(addr net.Addr, conn *net.Conn, mac net.HardwareAddr) Client {
	var c = Client{
		addr: addr,
		conn: conn,
		mac:  mac,
	}
	_ = c.creatID()
	c.status = 000
	return c
}

func (c *Client) logIn() {
	c.status = 500
}

func (c *Client) logOut() {
	c.status = 200
}

func (c *Client) creatID() error {
	c.id = (string)(rune(rand.Int()))
	return nil
}

type CliContainer struct {
	CliContainer map[string]Client //id as key
}

func (receiver CliContainer) Find(key string) Client {
	return receiver.CliContainer[key]
}

func (receiver CliContainer) New(c Client) error {
	receiver.CliContainer[c.id] = c
	return nil
}

func DiscoverClient(container *CliContainer, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.IPv4(0, 0, 0, 0),
	}

	// 监听 UDP 广播
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	fmt.Println("Listening for UDP broadcast on port", addr.Port)

	buf := make([]byte, 1024)

	for {
		// 读取 UDP 数据包
		n, srcAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}
		fmt.Printf("Received UDP broadcast from %s: %s\n", srcAddr, string(buf[:n]))
	}
}
