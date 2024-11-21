package WebAPI

import (
	"ConfigServer/APIGateway/APIWorkers"
	"ConfigServer/APIGateway/ClientAPICallers"
	"ConfigServer/clientManage"
	"encoding/json"
	"net"
	"net/http"
)

type functionResponse struct {
	FName      string `json:"f_name"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
	HttpStatus int    `json:"-"`
}

func sendFunctionResponse(rsp functionResponse, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(rsp.HttpStatus)
	return json.NewEncoder(w).Encode(rsp)
}

func function(w http.ResponseWriter, r *http.Request, body *[]byte) {

	fName, err := getFName(body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch fName {
	case "update_host_name":
		type updateHostNameReq struct {
			DestSysyncID []string `json:"dest_sysync_id"`
			FName        string   `json:"f_name"`
			DestIP       []string `json:"dest_ip"`
			DestPort     int      `json:"dest_port"`
			HostIP       string   `json:"host_ip"`
			HostPort     int      `json:"host_port"`
		}

		var requestData updateHostNameReq
		if err = json.Unmarshal(*body, &requestData); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		adders := make([]net.UDPAddr, 0, len(requestData.DestSysyncID)+len(requestData.DestIP))

		var destPort int
		if requestData.DestPort == 0 { // got null in json
			destPort = clientManage.UdpClientPort
		} else {
			destPort = requestData.DestPort
		}

		for _, v := range requestData.DestIP {

			if k := net.ParseIP(v); k != nil {
				c := net.UDPAddr{
					Port: destPort,
					IP:   k,
				}
				c.Port = destPort
				adders = adders[:len(adders)+1]
				adders[len(adders)-1] = c
			}
		}

		if requestData.HostIP == "" {
			requestData.HostIP = clientManage.CliUdpApiGateway.IP()
		}

		if len(adders) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sender := ClientAPICallers.NewUpdateHostName(
			adders,
			"update_host_name",
			requestData.HostIP,
			clientManage.CliUdpApiGateway.Port(),
		)

		nameServer := APIWorkers.HostNameReq{}
		nameServer.SetKeyWord("host_name_req")

		//t := clientManage.Schedule{
		//	ExecTime: time.Time{},
		//	Do: func() error {
		//		err := sender.Run()
		//		return err
		//	},
		//}
		_ = sender.Start()

		addErr := clientManage.CliUdpApiGateway.Add(&nameServer)
		if addErr == nil {
			_ = nameServer.Start()
		}
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusOK,
			Message:    "success",
			FName:      fName,
		}, w)

		//t2 := clientManage.Schedule{
		//	ExecTime: time.Time{},
		//	Do: func() error {
		//		err := sender.Run()
		//		return err
		//	},
		//}
		/*
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
		*/
	}

}
