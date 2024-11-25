package APIGateway

import (
	"ConfigServer/utils"
	"fmt"
	"net"
)

var CliUdpApiGateway *UDPAPIGateway = nil // definition of tcp port

type UDPMessage struct {
	Addr net.UDPAddr
	Text map[string]interface{}
}

type UDPAPIPortTemp struct {
	keyWord    string
	MessageQue *utils.Queue
	Gateway    *UDPAPIGateway

	EndRun chan bool
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
		_ = u.run()
	}()
	return nil
}

func (u *UDPAPIPortTemp) Stop() error {
	u.EndRun <- true
	return nil
}

func (u *UDPAPIPortTemp) NewMess(mess UDPMessage) {
	u.MessageQue.Append(mess)
}

// Init automatically called by gateway when it is added into
func (u *UDPAPIPortTemp) Init(gateway *UDPAPIGateway) {
	u.Gateway = gateway
	u.MessageQue = utils.NewQueue()
}

func (u *UDPAPIPortTemp) run() error { // a template to write APIs' definition
	stop := false

	for stop == false {
		reqPack := u.MessageQue.Pop().(UDPMessage)

		select {
		case <-u.EndRun:
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
	portList     map[string]UDPAPIPort // pointer point to a real port structure
	port         int
	ip           string
	statCode     int
	netInterface net.Interface

	udpListener *net.UDPConn

	endRun    chan bool
	initiated bool
}

func NewUDPAPIGateway(port int, ip string) *UDPAPIGateway {
	return &UDPAPIGateway{
		portList: make(map[string]UDPAPIPort),
		port:     port,
		ip:       ip,
		statCode: 000,
	}
}

func (a *UDPAPIGateway) Port() int {
	return a.port
}

func (a *UDPAPIGateway) IP() string {
	return a.ip
}

func (a *UDPAPIGateway) getIPInfo(interfaceName string) ([]net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 遍历所有接口，查找匹配的接口
	for _, iface := range interfaces {
		if iface.Name == interfaceName {
			// 获取该接口的地址
			addrs, err := iface.Addrs()
			if err != nil {
				return nil, err
			}

			// 收集有效的 IP 地址
			var ips []net.IP
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok {
					ips = append(ips, ipNet.IP)
				}
			}
			return ips, nil
		}
	}

	return nil, fmt.Errorf("interface %s not found", interfaceName)
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
		Port: a.port,
		IP:   net.ParseIP(a.ip),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	a.udpListener = conn
	a.initiated = true
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
	run() error
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

type ApiResponse struct {
	Fname   string `json:"f_name"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
