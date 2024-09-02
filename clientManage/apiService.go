//backend api service

package clientManage

import (
	"ConfigServer/utils"
	"bufio"
	"fmt"
	"net"
	"os"
)

type APIPort struct {
	keyWord string
	do      func()
}

type APIGateway struct { // listen api calls on a specific port
	portList map[string]*APIPort
	Port     string
	statCode int
	Protocol string // tcp or udp
}

func NewAPIGateway(port string, protocol string) *APIGateway {
	return &APIGateway{
		portList: make(map[string]*APIPort),
		Port:     port,
		Protocol: protocol,
		statCode: 000,
	}
}

func (a *APIGateway) Run() error {
	return nil
}

func (a *APIGateway) Stop() error {
	return nil
}

func (a *APIGateway) Add(port *APIPort) error {
	if a.portList == nil {
		a.portList = make(map[string]*APIPort)
	}
	if _, ok := a.portList[port.keyWord]; ok {
		return fmt.Errorf("port %s already exists", port.keyWord)
	} else {
		a.portList[port.keyWord] = port
		return nil
	}
}

func (a *APIGateway) Remove(keyWord string) error {
	if a.portList == nil {
		a.portList = make(map[string]*APIPort)
	}
	if _, ok := a.portList[keyWord]; ok {
		delete(a.portList, keyWord)
	} else {
		return fmt.Errorf("port %s not exists", keyWord)
	}
	return nil
}

func APIService(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {
			return
		}
	}(ln)

	fmt.Println("TCP listening on 6004...")

	for {
		// 接受客户端连接
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}

		// 处理客户端连接
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	//message, err := readUntilEndMarker(conn)
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		return
	}

	messJson, jsonErr := utils.JsonDecode([]byte(message))
	if jsonErr != nil {
		fmt.Println("Error unmarshalling JSON:", jsonErr.Error())
	}

	switch messJson["f_name"] {
	case "connect_resident_socket":
		fmt.Println("Connect req from client")

	case "disconnect_resident_socket":
		fmt.Println("Disconnect req from client")

	case "host_name":
		fmt.Println("Request host name")

	case "net":
	case "collect_input":
	case "reg":
	case "setting":
	}

	fmt.Printf("Received from client: %s\n", messJson)

	_, sendErr := conn.Write([]byte("{\"Success\": 1}\n"))
	if sendErr != nil {
		fmt.Println("Error sending response:", sendErr.Error())
	}
	return
}
