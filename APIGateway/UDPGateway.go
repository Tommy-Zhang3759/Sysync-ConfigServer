package APIGateway

import (
	"ConfigServer/utils"
	"fmt"
	"net"
)

type UDPMessage struct {
	Addr net.UDPAddr
	Text map[string]interface{}
}

type UDPAPIPortTemp struct {
	keyWord    string
	messageQue *utils.Queue
	Gateway    *UDPAPIGateway

	endRun chan bool
}

func (u *UDPAPIPortTemp) SetKeyWord(key string) {
	u.keyWord = key
	return
}

func (u *UDPAPIPortTemp) GetKeyWord() string {
	return u.keyWord
}

func (u *UDPAPIPortTemp) Start() error {
	go func() {
		_ = u.Run()
	}()
	return nil
}

func (u *UDPAPIPortTemp) Stop() error {
	u.endRun <- true
	return nil
}

func (u *UDPAPIPortTemp) NewMess(mess UDPMessage) {
	u.messageQue.Append(mess)
}

func (u *UDPAPIPortTemp) Init(gateway *UDPAPIGateway) {
	u.Gateway = gateway
	u.messageQue = utils.NewQueue()
}

func (u *UDPAPIPortTemp) Run() error { // a template to write APIs' definition
	stop := false

	for stop == false {
		reqPack := u.messageQue.Pop().(UDPMessage)

		select {
		case <-u.endRun:
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:
			err := u.Gateway.SendMess([]byte(reqPack.Text["f_name"].(string)), reqPack.Addr)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

type UDPAPIGateway struct { // listen api calls on a specific port
	portList map[string]UDPAPIPort // pointer point to a real port structure
	Port     int
	statCode int

	udpListener *net.UDPConn

	endRun chan bool
	inited bool
}

func (a *UDPAPIGateway) SendMess(mess []byte, destIPs ...net.UDPAddr) error {
	for _, destIP := range destIPs {
		conn, err := net.DialUDP("udp", nil, &destIP)
		if err != nil {
			fmt.Println("Error connecting:", err)
			return err
		}

		_, err = conn.Write(mess)
		if err != nil {
			fmt.Println("Error sending UDP message:", err)
			return err
		}
		fmt.Println("Message sent to", destIP.String(), ": ", string(mess))
		_ = conn.Close()
	}
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
	a.inited = true
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

				if b := messJson["f_name"] == nil; b {
					_ = a.SendMess([]byte("{\"error\":\"invalid key\"}"))
				}

				a.portList[messJson["f_name"].(string)].NewMess(UDPMessage{
					Addr: addr,
					Text: messJson,
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

type UDPAPIPort interface {
	Run() error
	Start() error
	Stop() error
	NewMess(mess UDPMessage)
	GetKeyWord() string
	Init(gateway *UDPAPIGateway)
	SetKeyWord(key string)
}

func (a *UDPAPIGateway) Add(port UDPAPIPort) error {
	if _, ok := a.portList[port.GetKeyWord()]; ok {
		return fmt.Errorf("port %s already exists", port.GetKeyWord())
	} else {
		port.Init(a)
		a.portList[port.GetKeyWord()] = port
		return nil
	}
}

func (a *UDPAPIGateway) Remove(keyWord string) error {
	if _, ok := a.portList[keyWord]; ok {
		delete(a.portList, keyWord)
	} else {
		return fmt.Errorf("port %s not exists", keyWord)
	}
	return nil
}
