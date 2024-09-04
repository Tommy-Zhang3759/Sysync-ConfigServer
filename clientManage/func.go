package clientManage

import (
	"ConfigServer/APIGateway"
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var UdpHostPort = 6004

var CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpHostPort) // definition of tcp port

type Call interface {
	Run() error
	Stop() error
}

type TCPCallMessage struct {
	targetClient *list.List
	body         map[string]interface{}
}

func (c *TCPCallMessage) Run() error {

	return nil
}

func (c *TCPCallMessage) Stop() error {
	return nil
}

type UDPCallMessage struct {
	TargetIP *list.List
	Body     map[string]interface{}
	Port     string
	stop     bool
}

func (c *UDPCallMessage) Run() error {

	var wg sync.WaitGroup

	// 遍历链表中的所有IP地址
	for e := c.TargetIP.Front(); e != nil; e = e.Next() {
		targetIP, ok := e.Value.(string)
		if !ok {
			fmt.Println("Error: Invalid IP address in list")
			continue
		}

		wg.Add(1)
		go func() { //message sender
			err := func(ip string) error {
				defer wg.Done()

				serverAddr, err := net.ResolveUDPAddr("udp", ip+":"+c.Port)
				if err != nil {

					fmt.Println("Error resolving UDP address:", err)
					return err
				}

				conn, err := net.DialUDP("udp", nil, serverAddr)
				if err != nil {
					fmt.Println("Error connecting:", err)
					return err
				}
				defer func(conn *net.UDPConn) {
					err := conn.Close()
					if err != nil {
						return
					}
				}(conn)

				_, err = conn.Write(c.bodyJson())
				if err != nil {
					fmt.Println("Error sending UDP message:", err)
					return err
				}
				fmt.Println("Message sent to", serverAddr.String())
				fmt.Println("Body:", string(c.bodyJson()))
				return nil
			}(targetIP)
			if err != nil {
				return //handle err in FUTURE
			}
		}()
	}

	wg.Wait() // 等待所有 Goroutine 完成
	return nil
}

func (c *UDPCallMessage) Cancel() error {
	c.stop = true
	return nil
}

func (c *UDPCallMessage) bodyJson() []byte {
	mess, _ := json.Marshal(c.Body)
	return mess
}

func NewUDPCallMess() *UDPCallMessage {
	return &UDPCallMessage{
		TargetIP: list.New(),
		Body: map[string]interface{}{
			"f_name": "undefined",
		},
	}
}

func sendCommand(cmd string) {

}

func HostNameRequester() Schedule {
	var api = APIGateway.UDPAPIPort{
		KeyWord: "updateHostName",
		Do: func(req map[string]interface{}, addr net.UDPAddr) error {

			type ClientRequest struct {
				Mac string `json:"mac"`
			}

			type HostInfo struct {
				IP       string `json:"ip"`
				HostName string `json:"host_name"`
			}

			macMap := func() map[string]HostInfo {
				jsonFile, err := os.Open("resources/mac_ip_host_name.json")
				if err != nil {
					log.Fatalf("Can not open file: %v", err)
				}
				defer func(jsonFile *os.File) {
					_ = jsonFile.Close()
				}(jsonFile)

				// 读取文件内容
				byteValue, err := ioutil.ReadAll(jsonFile)
				if err != nil {
					log.Fatalf("无法读取文件内容: %v", err)
				}

				// 定义一个map来存储JSON数据
				macMap := make(map[string]HostInfo)

				// 解析JSON到map中
				if err := json.Unmarshal(byteValue, &macMap); err != nil {
					log.Fatalf("JSON解析失败: %v", err)
				}
				return macMap
			}()

			cliInfo := (macMap)[req["mac"].(string)]

			rspMessage, err := json.Marshal(cliInfo)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return err
			}

			sendAddr := net.UDPAddr{
				Port: 0,
				IP:   net.ParseIP("0.0.0.0"),
			}
			conn, err := net.ListenUDP("udp", &sendAddr)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			defer func(conn *net.UDPConn) {
				_ = conn.Close()
			}(conn)

			// 目标地址
			targetAddr := net.UDPAddr{
				IP:   net.ParseIP(addr.String()),
				Port: addr.Port,
			}

			// 发送数据
			_, err = conn.WriteTo(rspMessage, &targetAddr)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return err
			}

			return nil
		},
	}

	return Schedule{
		execTime: time.Time{},
		do: func() error {
			err := CliUdpApiGateway.Add(&api)
			if err != nil {
				return err
			}
			return nil
		},
	}

}

func getInput() {

}

func getOutput() {

}

func modifyREG(key string, subkey string) {

}

func collectREG() {

}
