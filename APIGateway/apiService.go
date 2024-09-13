//backend api service

package APIGateway

func NewUDPAPIGateway(port int) *UDPAPIGateway {
	return &UDPAPIGateway{
		portList: make(map[string]*UDPAPIPort),
		Port:     port,
		statCode: 000,
	}
}

//func handleConnection(conn net.Conn) {
//	defer func(conn net.Conn) {
//		err := conn.Close()
//		if err != nil {
//			return
//		}
//	}(conn)
//
//	//message, err := readUntilEndMarker(conn)
//	message, err := bufio.NewReader(conn).ReadString('\n')
//	if err != nil {
//		fmt.Println("Error reading from connection:", err.Error())
//		return
//	}
//
//	messJson, jsonErr := utils.JsonDecode([]byte(message))
//	if jsonErr != nil {
//		fmt.Println("Error unmarshalling JSON:", jsonErr.Error())
//	}
//
//	switch messJson["f_name"] {
//	case "connect_resident_socket":
//		fmt.Println("Connect req from client")
//
//	case "disconnect_resident_socket":
//		fmt.Println("Disconnect req from client")
//
//	case "host_name":
//		fmt.Println("Request host name")
//
//	case "net":
//	case "collect_input":
//	case "reg":
//	case "setting":
//	}
//
//	fmt.Printf("Received from client: %s\n", messJson)
//
//	_, sendErr := conn.Write([]byte("{\"Success\": 1}\n"))
//	if sendErr != nil {
//		fmt.Println("Error sending response:", sendErr.Error())
//	}
//	return
//}
