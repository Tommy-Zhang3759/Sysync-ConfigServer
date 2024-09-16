package clientManage

import (
	"ConfigServer/APIGateway"
)

var UdpHostPort = 6004
var UdpClientPort = 6003

var CliUdpApiGateway = APIGateway.NewUDPAPIGateway(UdpHostPort) // definition of tcp port

func SendCommand2Host(cmd string) {

}

//func HostNameRequester() Schedule {
//
//	return Schedule{
//		ExecTime: time.Time{},
//		Do: func() error {
//			err := CliUdpApiGateway.Add(&api)
//			if err != nil {
//				return err
//			}
//			return nil
//		},
//	}
//
//}

func getInput() {

}

func getOutput() {

}

func modifyREG(key string, subkey string) {

}

func collectREG() {

}
