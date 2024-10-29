package WebAPI

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"
	"ConfigServer/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request, q url.Values) {
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
			case "cliInfo":
				cliInfo(w, r)
			default:

				http.Error(w, "API not found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Invalid API request", http.StatusBadRequest)
		}
	} else if r.Method == http.MethodGet {
		if len(pathSegments) > 0 {
			switch pathSegments[0] {
			case "cliInfo":
				cliInfo(w, r)
			default:
				http.Error(w, "API not found", http.StatusNotFound)
			}
		} else {
			http.Error(w, "Invalid API request", http.StatusBadRequest)
		}
	}
}

func cliInfo(w http.ResponseWriter, r *http.Request) {
	/*
	* return a list containing all host names if no arguments are included
	* otherwise, return detailed information
	 */

	var responseData []byte
	var err error

	idList := r.URL.Query()["id"]

	if idList != nil && len(idList) > 0 {
		oriInfo, _, _ := clientManage.DetailedInfo(idList)

		var friendly []clientManage.FriendlyClient

		for i := range oriInfo {
			f, _ := oriInfo[i].HumanFriendly()
			friendly = append(friendly, f)
		}

		responseData, err = json.Marshal(friendly)

	} else {
		responseData, err = json.Marshal(clientManage.AllHostName())
	}

	if err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseData)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
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

		nameServer := APIGateway.HostNameReq{}
		nameServer.SetKeyWord("host_name_req")

		//t := clientManage.Schedule{
		//	ExecTime: time.Time{},
		//	Do: func() error {
		//		err := sender.Run()
		//		return err
		//	},
		//}
		_ = sender.Run()

		addErr := clientManage.CliUdpApiGateway.Add(&nameServer)
		if addErr == nil {
			_ = nameServer.Run()
		}

		//t2 := clientManage.Schedule{
		//	ExecTime: time.Time{},
		//	Do: func() error {
		//		err := sender.Run()
		//		return err
		//	},
		//}

	case "send_command_to_host":

		sender := APIGateway.MessSending{
			Dest: addrs,
			MessContent: map[string]interface{}{
				"f_name":  "run_command",
				"command": bodyJson["command"].(string),
			},
		}

		//commandServer := APIGateway.CommandReq{}
		//commandServer.SetKeyWord("command_req")
		//
		_ = sender.Run()
		//
		//addErr := clientManage.CliUdpApiGateway.Add(&commandServer)
		//if addErr == nil {
		//	_ = commandServer.Run()
		//}
		//
		//
	case "set_server_info":
		sender := APIGateway.MessSending{
			Dest: addrs,
			MessContent: map[string]interface{}{
				"f_name":      "set_server_info",
				"server_ip":   bodyJson["server_ip"].(string),
				"server_port": bodyJson["server_port"].(int),
			},
		}
		_ = sender.Run()
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
