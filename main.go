package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// 监听端口 8080
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer func(ln net.Listener) {
		err := ln.Close()
		if err != nil {

		}
	}(ln)

	fmt.Println("TCP 服务器正在监听端口 8080...")

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

		}
	}(conn)

	// 读取客户端发送的数据
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("收到客户端消息: %s", message)

	// 向客户端发送响应
	_, _ = conn.Write([]byte("消息已收到\n"))
}
