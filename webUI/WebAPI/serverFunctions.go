package WebAPI

import (
	"ConfigServer/APIGateway"
	"ConfigServer/APIGateway/APIWorkers"
	"ConfigServer/APIGateway/ClientAPICallers"
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

type functionResponse struct {
	FName      string `json:"f_name"`
	Message    string `json:"message,omitempty"`
	Error      string `json:"error,omitempty"`
	HttpStatus int    `json:"-"`
}

type webReqTemp struct {
	DestSysyncID []string `json:"dest_sysync_id,omitempty"`
	FName        string   `json:"f_name"`
	DestIP       []string `json:"dest_ip,omitempty"`
	DestPort     int      `json:"dest_port,omitempty"`
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
		updateHostName(w, body, fName)
	case "run_command":
		runCommand(w, body, fName)
	}

}
func formatDestAddr(destIPs *[]string, destPort int) ([]net.UDPAddr, error) {
	adders := make([]net.UDPAddr, 0, len(*destIPs))

	if len(*destIPs) == 1 && (*destIPs)[0] == "" {
		return nil, nil
	}

	if destPort <= 0 && len(*destIPs) > 1 && (*destIPs)[0] != "" {
		return nil, errors.New("invalid port number")
	} else {
		destPort = APIGateway.CliUdpApiGateway.Port()
	}

	for _, v := range *destIPs {

		if k := net.ParseIP(v); k != nil {
			c := net.UDPAddr{
				Port: destPort,
				IP:   k,
			}
			c.Port = destPort
			adders = adders[:len(adders)+1]
			adders[len(adders)-1] = c
		} else {
			return nil, errors.New("invalid ip address")
		}
	}
	return adders, nil
}
func updateHostName(w http.ResponseWriter, body *[]byte, fName string) {
	type reqTemp struct {
		webReqTemp
		HostIP   string `json:"host_ip,omitempty"`
		HostPort int    `json:"host_port,omitempty"`
	}

	var requestData reqTemp
	if err := json.Unmarshal(*body, &requestData); err != nil {
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      err.Error(),
		}, w)
		return
	}

	if requestData.HostIP == "" {
		requestData.HostIP = APIGateway.CliUdpApiGateway.IP()
	}
	if requestData.HostPort == 0 {
		requestData.HostPort = APIGateway.CliUdpApiGateway.Port()
	}
	sender := ClientAPICallers.NewUpdateHostName(
		requestData.HostIP,
		requestData.HostPort,
	)

	address, err := formatDestAddr(&requestData.DestIP, requestData.DestPort)
	if err != nil {
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      err.Error(),
		}, w)
		return
	}

	if err = sender.MoreDestBySysyncID(requestData.DestSysyncID...); err != nil {
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      err.Error(),
		}, w)
		return
	}
	sender.MoreDestByIP(address...)

	namingService := APIWorkers.HostNameReq{}
	namingService.SetKeyWord(fName)

	//t := clientManage.Schedule{
	//	ExecTime: time.Time{},
	//	Do: func() error {
	//		err := sender.Run()
	//		return err
	//	},
	//}

	addErr := APIGateway.CliUdpApiGateway.Add(&namingService)
	if addErr == nil {
		_ = namingService.Start()
	}

	_ = sender.Run()

	_ = sendFunctionResponse(functionResponse{
		HttpStatus: http.StatusOK,
		Message:    "success",
		FName:      fName,
	}, w)
}

func runCommand(w http.ResponseWriter, body *[]byte, fName string) {
	type reqTemp struct {
		webReqTemp
		Command string `json:"command,omitempty"`
	}

	var requestData reqTemp
	if err := json.Unmarshal(*body, &requestData); err != nil {
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      err.Error(),
		}, w)
		return
	}

	sender := ClientAPICallers.NewRunCmd(requestData.Command)

	address, err := formatDestAddr(&requestData.DestIP, requestData.DestPort)
	if err != nil {
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      err.Error(),
		}, w)
		return
	}

	if err = sender.MoreDestBySysyncID(requestData.DestSysyncID...); err != nil {
		_ = sendFunctionResponse(functionResponse{
			HttpStatus: http.StatusBadRequest,
			Error:      err.Error(),
		}, w)
		return
	}
	sender.MoreDestByIP(address...)

	namingService := APIWorkers.HostNameReq{}
	namingService.SetKeyWord(fName)

	//t := clientManage.Schedule{
	//	ExecTime: time.Time{},
	//	Do: func() error {
	//		err := sender.Run()
	//		return err
	//	},
	//}

	addErr := APIGateway.CliUdpApiGateway.Add(&namingService)
	if addErr == nil {
		_ = namingService.Start()
	}

	_ = sender.Run()

	_ = sendFunctionResponse(functionResponse{
		HttpStatus: http.StatusOK,
		Message:    "success",
		FName:      fName,
	}, w)
}
