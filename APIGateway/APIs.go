package APIGateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type HostNameReq struct {
	UDPAPIPortTemp
}

func (u *HostNameReq) Run() error {
	stop := false

	for stop == false {
		reqPack := u.messageQue.Pop().(UDPMessage)

		select {
		case <-u.endRun:
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:

			type ClientRequest struct {
				Mac string `json:"mac"`
			}

			type HostInfo struct { // use database in FUTURE
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

			cliInfo := macMap[reqPack.Text["mac"].(string)]
			type temp struct {
				FName string `json:"f_name"`
				Host  string `json:"host_name"`
				IP    string `json:"host_ip"`
			}

			rsp := temp{
				FName: "host_name_offer",
				Host:  cliInfo.HostName,
				IP:    cliInfo.IP,
			}

			rspJson, err := json.Marshal(rsp)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return err
			}

			err = u.Gateway.SendMess(rspJson, reqPack.Addr)
			if err != nil {
				fmt.Println("Error sending message:", err)
				return err
			}

		}
	}
	return nil
}

type sendCommandToHost struct {
	UDPAPIPortTemp
}

type MessSending struct {
	UDPAPIPortTemp
	Dest        []net.UDPAddr
	MessContent map[string]interface{}
}

func (m *MessSending) name() {

}
func (m *MessSending) Run() error {
	err := m.Gateway.SendMess(m.bodyJson(), m.Dest...)
	return err
}

func (m *MessSending) bodyJson() []byte {
	mess, _ := json.Marshal(m.MessContent)
	return mess
}
