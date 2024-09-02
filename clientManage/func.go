package clientManage

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
)

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
	Server   func()
	stop     bool
}

func (c *UDPCallMessage) Run() error {
	c.Server() // need to wait until server start successfully in FUTURE

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
				fmt.Println("Message sent to", ip)
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

func UpdateHostName(t *UDPCallMessage, destStrings []string, hostIP string, hostPort string) {
	for i := range destStrings {
		t.TargetIP.PushBack((destStrings)[i])
	}

	t.Body["f_name"] = "updateHostName"
	t.Body["host_ip"] = hostIP
	t.Body["host_port"] = hostPort
	t.Port = hostPort // use the same port as the sending service to save port usage

	t.Server = func() {
		serverAddr := t.Body["host_ip"].(string) + ":" + t.Body["host_port"].(string)
		listener, err := net.Listen("tcp", serverAddr)

		if err != nil {
			fmt.Println("Error creating server socket:", err)
			return
		}

		fmt.Printf("Host name update service listening on %s\n", serverAddr)

		type ClientRequest struct {
			Mac string `json:"mac"`
		}

		//------------------------------tmp use----------------------------------------

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

		//------------------------------tmp use----------------------------------------
		go func() {
			for !t.stop {
				// 接收客户端连接
				clientConn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err)
					continue
				}

				// 处理客户端请求
				go func(conn net.Conn, macMap *map[string]HostInfo) {
					defer func(conn net.Conn) {
						_ = conn.Close()
					}(conn)

					// 接收客户端的请求消息
					buffer := make([]byte, 1024)
					n, err := conn.Read(buffer)
					if err != nil {
						fmt.Println("Error reading from connection:", err)
						return
					}

					// 将请求消息解析为字典
					var req ClientRequest
					err = json.Unmarshal(buffer[:n], &req)
					if err != nil {
						fmt.Println("Error unmarshalling JSON:", err)
						return
					}

					fmt.Printf("Received from client: %+v\n", req)

					// 模拟从 DataFrame 获取客户端信息
					cliInfo := (*macMap)[req.Mac]

					// 将响应消息转换为 JSON
					responseMessage, err := json.Marshal(cliInfo)
					if err != nil {
						fmt.Println("Error marshalling JSON:", err)
						return
					}

					// 发送响应消息回客户端
					_, err = conn.Write(responseMessage)
					if err != nil {
						fmt.Println("Error sending response:", err)
						return
					}

				}(clientConn, &macMap)
			}
			defer func(listener net.Listener) {
				_ = listener.Close()
			}(listener)
		}()
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
