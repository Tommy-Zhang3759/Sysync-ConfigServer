package clientManage

import (
	"ConfigServer/APIGateway"
	"net"
	"time"
)

var UdpHostPort = 6004
var UdpClientPort = 6003

var CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpHostPort) // definition of tcp port

func SendCommand2Host(cmd string) {

}

func HostNameRequester() Schedule {
	var api = APIGateway.UDPAPIPortTemp{
		keyWord: "host_name_req",
		Do: func(req map[string]interface{}, addr net.UDPAddr) error {

		},
	}

	return Schedule{
		execTime: time.Time{},
		do: func() error {
			err := CliUdpApiGateway.Add(&api)
			if err != nil {
				return err
			}
			return nil
		},
	}

}

func getInput() {

}

func getOutput() {

}

func modifyREG(key string, subkey string) {

}

func collectREG() {

}
