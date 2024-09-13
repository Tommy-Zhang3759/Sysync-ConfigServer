package APIGateway

import (
	"ConfigServer/utils"
	"fmt"
	"net"
)

type reqMessage struct {
	Source net.UDPAddr
	Text   map[string]interface{}
}

type UDPAPIPort struct {
	KeyWord    string
	messageQue utils.Queue
	Gateway    *UDPAPIGateway

	endRun chan bool
}

func (u *UDPAPIPort) newMess(mess reqMessage) {
	u.messageQue.Append(mess)
}

func (u *UDPAPIPort) run() error { // a template to write APIs' definition
	stop := false

	for stop == false {
		reqPack := u.messageQue.Pop().(reqMessage)

		select {
		case <-u.endRun:
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:
			err := u.Gateway.sendMess(reqPack.Source, []byte(reqPack.Text["f_name"].(string)))
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

type UDPAPIGateway struct { // listen api calls on a specific port
	portList map[string]*UDPAPIPort
	Port     int
	statCode int

	udpListener *net.UDPConn

	endRun chan bool
}

func (a *UDPAPIGateway) sendMess(destIP net.UDPAddr, mess []byte) error {

	conn, err := net.DialUDP("udp", nil, &destIP)
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

	_, err = conn.Write(mess)
	if err != nil {
		fmt.Println("Error sending UDP message:", err)
		return err
	}
	fmt.Println("Message sent to", destIP.String(), ": ", mess)
	return nil

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
	a.udpListener = conn
	return nil
}

func (a *UDPAPIGateway) Run() error {
	var stop = false
	buffer := make([]byte, 1024)
	for stop == false {
		n, addr, err := a.udpListener.ReadFrom(buffer)

		select {
		case <-a.endRun:
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:
			if err != nil {
				fmt.Println("Error reading form UDP API message:", err.Error())
				continue
			}

			go func(buffer []byte, addr net.UDPAddr) {
				messJson, jsonErr := utils.JsonDecode(buffer)
				if jsonErr != nil {
					fmt.Println("Error decoding form UDP API message:", jsonErr.Error())
				}

				a.portList[messJson["f_name"].(string)].newMess(reqMessage{
					Source: addr,
					Text:   messJson,
				})
			}(buffer[:n], *addr.(*net.UDPAddr))
		}

		//fmt.Printf("Received %s from %s\n", string(buffer[:n]), addr)

	}
	return nil
}

func (a *UDPAPIGateway) Stop() error {
	_ = a.udpListener.Close()
	a.endRun <- true

	return nil
}

func (a *UDPAPIGateway) Add(port *UDPAPIPort) error {
	if a.portList == nil {
		a.portList = make(map[string]*UDPAPIPort)
	}

	if _, ok := a.portList[port.KeyWord]; ok {
		return fmt.Errorf("port %s already exists", port.KeyWord)
	} else {
		port.Gateway = a
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
