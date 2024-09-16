package webUI

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"
	"ConfigServer/utils"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	//_, err := fmt.Fprintf(w, "Hello from server on port 8080! You've requested: %s\n", r.URL.Path)
	//if err != nil {
	//	return
	//}

	reqPath := r.URL.Path
	query := r.URL.Query()

	if strings.HasPrefix(reqPath, "/api/") {
		handleAPI(w, r, query)
		return
	}

	// 设置静态文件的根目录
	staticDir := "./webUI/static/"

	// 构造文件的完整路径
	path := filepath.Join(staticDir, reqPath)

	// 返回本地文件
	http.ServeFile(w, r, path)
}

func handleAPI(w http.ResponseWriter, r *http.Request, q url.Values) {
	apiPath := strings.TrimPrefix(r.URL.Path, "/api/")
	pathSegments := strings.Split(apiPath, "/")

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(r.Body)

		_, err = fmt.Printf("POST request body: %s\n", body)

		if len(pathSegments) > 0 {
			switch pathSegments[0] {
			case "login":

				adminLogin(w, r, q.Get("username"), q.Get("password"))
			case "logout":

				adminLogout(w, r, q.Get("username"))
			case "cfg":
				command(w, r, &body)
			case "func":
				function(w, r, &body)
			default:

				http.Error(w, "API not found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Invalid API request", http.StatusBadRequest)
		}
	}
}

func function(w http.ResponseWriter, r *http.Request, body *[]byte) {
	bodyJson, _ := utils.JsonDecode(*body)

	var destStrings []string
	if dest, ok := bodyJson["dest_ip"].([]interface{}); ok {
		for _, v := range dest {
			if str, ok := v.(string); ok {
				destStrings = append(destStrings, str)
			}
		}
	} else {
		destStrings = append(destStrings, "0.0.0.0")
	}

	var addrs []net.UDPAddr
	for _, ip := range destStrings {
		addr, err := utils.ParseUDPAddr(ip, bodyJson["dest_port"].(string))
		if err != nil {
			fmt.Println("Error parsing address:", err)
			continue
		}
		addrs = append(addrs, addr)
	}

	switch bodyJson["f_name"] {
	case "update_host_name":

		sender := APIGateway.MessSending{
			Dest: addrs,
			MessContent: map[string]interface{}{
				"f_name":    "update_host_name",
				"host_ip":   bodyJson["host_ip"].(string),
				"host_port": clientManage.UdpHostPort,
				//"host_port": bodyJson["host_port"].(string),
			},
		}

		//t := clientManage.Schedule{
		//	ExecTime: time.Time{},
		//	Do: func() error {
		//		err := sender.Run()
		//		return err
		//	},
		//}
		_ = clientManage.CliUdpApiGateway.Add(&sender)
		_ = sender.Run()

		//t2 := clientManage.Schedule{
		//	ExecTime: time.Time{},
		//	Do: func() error {
		//		err := sender.Run()
		//		return err
		//	},
		//}

	case "send_command_to_host":
		//t := clientManage.NewUDPCallMess()
		//
		//var destStrings []string
		//if dest, ok := bodyJson["dest_ip"].([]interface{}); ok {
		//	for _, v := range dest {
		//		if str, ok := v.(string); ok {
		//			destStrings = append(destStrings, str)
		//		}
		//	}
		//}
		//for _, v := range destStrings {
		//	t.TargetIP.PushBack(v)
		//}
		//
		//if bodyJson["dest_port"].(string) != "" {
		//	t.Port = bodyJson["dest_port"].(string)
		//} else {
		//	t.Port = (string)(rune(clientManage.UdpClientPort))
		//}
		//
		//t.Body["f_name"] = "run_command"
		//// t.Body["host_ip"] = bodyJson["host_ip"].(string)
		//// t.Body["host_port"] = clientManage.UdpHostPort
		//t.Body["command"] = bodyJson["command"]
		//
		//err := t.Run()
		//if err != nil {
		//	println(err.Error())
		//}
		// add t into task list and return an tID
	}

}

func command(w http.ResponseWriter, r *http.Request, body *[]byte) {
	print(body)
	//bodyJson, _ := utils.JsonDecode(*body)
	//com := bodyJson["command"].(string)
	//console.Handler(com)

}

func adminLogin(w http.ResponseWriter, r *http.Request, userName string, password string) {
	if userName == "" || password == "" {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
	}

}

func adminLogout(w http.ResponseWriter, r *http.Request, userName string) {
	if userName == "" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
	}

}

func StartServer(port string, handlerFunc http.HandlerFunc) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	clientManage.CliUdpApiGateway.Init()
	go clientManage.CliUdpApiGateway.Run()

	fmt.Printf("Starting server at port %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server on port %s failed to start: %v\n", port, err)
	}
}
