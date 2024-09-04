package APIGateway

import (
	"ConfigServer/utils"
	"bufio"
	"fmt"
	"net"
)

type TCPAPIPort struct {
	KeyWord string
	Do      func(req interface{}, conn net.Conn) error
}

type TCPAPIGateway struct { // listen api calls on a specific port
	portList map[string]*TCPAPIPort
	Port     int
	statCode int

	tcpListener net.Listener

	endRun chan bool
}

func (a *TCPAPIGateway) Init() error {

	ln, err := net.Listen("tcp", ":"+string(rune(a.Port)))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	a.tcpListener = ln
	return nil
}

func (a *TCPAPIGateway) Run() error {
	var stop = false
	for stop == false {
		// 接受客户端连接
		conn, err := a.tcpListener.Accept()

		select {
		case <-a.endRun:
			// 接收到停止信号，退出goroutine
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:
			if err != nil {
				fmt.Println("Error accepting API connection:", err.Error())
				continue
			}

			// 处理客户端连接
			go func(conn net.Conn) {
				defer func(conn net.Conn) {
					err := conn.Close()
					if err != nil {
						fmt.Println("Error closing connection:", err.Error())
						return
					}
				}(conn)

				message, err := bufio.NewReader(conn).ReadString('\n')
				if err != nil {
					fmt.Println("Error reading from API call connection:", err.Error())
					return
				}

				messJson, jsonErr := utils.JsonDecode([]byte(message))
				if jsonErr != nil {
					fmt.Println("Error unmarshalling JSON from API call:", jsonErr.Error())
				}

				err = a.portList[messJson["f_name"].(string)].Do(messJson, conn)
				if err != nil {
					// write log
					return
				}

			}(conn)
		}
	}

	return nil
}

func (a *TCPAPIGateway) Stop() error {
	if a.tcpListener != nil {
		err := a.tcpListener.Close()
		if err != nil {
			fmt.Println("Error closing tcpListener:", err.Error())
			return err
		}
	}

	return nil
}

func (a *TCPAPIGateway) Add(port *TCPAPIPort) error {
	if a.portList == nil {
		a.portList = make(map[string]*TCPAPIPort)
	}
	if _, ok := a.portList[port.KeyWord]; ok {
		return fmt.Errorf("port %s already exists", port.KeyWord)
	} else {
		a.portList[port.KeyWord] = port
		return nil
	}
}

func (a *TCPAPIGateway) Remove(keyWord string) error {
	if a.portList == nil {
		a.portList = make(map[string]*TCPAPIPort)
	}
	if _, ok := a.portList[keyWord]; ok {
		delete(a.portList, keyWord)
	} else {
		return fmt.Errorf("port %s not exists", keyWord)
	}
	return nil
}

func (a *TCPAPIGateway) WriteTo(mess []byte) error {
	return nil
}
