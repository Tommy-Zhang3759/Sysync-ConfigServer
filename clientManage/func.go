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
	"time"
)

var UdpHostPort = 6004
var UdpClientPort = 6003

var CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpHostPort) // definition of tcp port

type UDPCallMessage struct {
	TargetIP *list.List
	Body     map[string]interface{}
	Port     string
	stop     bool
}

func (c *UDPCallMessage) Run() error {

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

func SendCommand2Host(cmd string) {

}

func HostNameRequester() Schedule {
	var api = APIGateway.UDPAPIPort{
		KeyWord: "host_name_req",
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

			cliInfo := macMap[req["mac"].(string)]
			type resTemp struct {
				FName string `json:"f_name"`
				Host  string `json:"host_name"`
				IP    string `json:"host_ip"`
			}

			rsp := resTemp{
				FName: "host_name_offer",
				Host:  cliInfo.HostName,
				IP:    cliInfo.IP,
			}

			rspJson, err := json.Marshal(rsp)
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
				IP:   addr.IP,
				Port: UdpClientPort,
			}

			// 发送数据
			_, err = conn.WriteTo(rspJson, &targetAddr)
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
