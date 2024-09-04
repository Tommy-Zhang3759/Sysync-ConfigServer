package APIGateway

import (
	"ConfigServer/utils"
	"fmt"
	"net"
)

type UDPAPIPort struct {
	KeyWord string
	Do      func(req map[string]interface{}, addr net.UDPAddr) error
}

type UDPAPIGateway struct { // listen api calls on a specific port
	portList map[string]*UDPAPIPort
	Port     int
	statCode int

	udpListener net.UDPConn

	endRun chan bool
}

func (a *UDPAPIGateway) Init() error {
	addr := net.UDPAddr{
		Port: a.Port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	a.udpListener = *conn
	return nil
}

func (a *UDPAPIGateway) Run() error {
	var stop = false
	buffer := make([]byte, 1024)
	for stop == false {
		n, addr, err := a.udpListener.ReadFrom(buffer)
		select {
		case <-a.endRun:
			// 接收到停止信号，退出goroutine
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:
			if err != nil {
				fmt.Println("Error reading form UDP API message:", err.Error())
				continue
			}
			rmtUDPAddr := *addr.(*net.UDPAddr)
			go func(buffer []byte, addr net.UDPAddr) {
				messJson, jsonErr := utils.JsonDecode(buffer)
				if jsonErr != nil {
					fmt.Println("Error decoding form UDP API message:", jsonErr.Error())
				}
				a.portList[messJson["f_name"].(string)].Do(messJson, addr)
			}(buffer[:n], rmtUDPAddr)
		}

		//fmt.Printf("Received %s from %s\n", string(buffer[:n]), addr)

	}
	return nil
}

func (a *UDPAPIGateway) Stop() error {

	err := a.udpListener.Close()
	if err != nil {
		fmt.Println("Error closing udpListener:", err.Error())
		return err
	}
	return nil
}

func (a *UDPAPIGateway) Add(port *UDPAPIPort) error {
	if a.portList == nil {
		a.portList = make(map[string]*UDPAPIPort)
	}
	if _, ok := a.portList[port.KeyWord]; ok {
		return fmt.Errorf("port %s already exists", port.KeyWord)
	} else {
		a.portList[port.KeyWord] = port
		return nil
	}
}

func (a *UDPAPIGateway) Remove(keyWord string) error {
	if a.portList == nil {
		a.portList = make(map[string]*UDPAPIPort)
	}
	if _, ok := a.portList[keyWord]; ok {
		delete(a.portList, keyWord)
	} else {
		return fmt.Errorf("port %s not exists", keyWord)
	}
	return nil
}
