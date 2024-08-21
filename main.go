package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func init() {

}

func main() {
	// 监听端口 8080
	ln, err := net.Listen("tcp", ":6004")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {

		}
	}(ln)

	fmt.Println("TCP 服务器正在监听端口 6004...")

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

func jsonDecode(jsonData []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}
	return result, err
}

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	// 读取客户端发送的数据
	message, _ := bufio.NewReader(conn).ReadString('\n')

	messJson, jsonErr := jsonDecode([]byte(message))
	if jsonErr != nil {
		fmt.Println("Error unmarshalling JSON:", jsonErr.Error())
		return
	}

	switch messJson["f_name"] {
	case "connect":
	}

	fmt.Printf("收到客户端消息: %s\n", message)

	// 向客户端发送响应
	_, sendErr := conn.Write([]byte("Success\n"))
	if sendErr != nil {
		fmt.Println("Error sending response")
	}
	return
}
