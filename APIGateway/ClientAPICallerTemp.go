package APIGateway

import (
	"ConfigServer/clientManage"
	"encoding/json"
	"net"
)

type ClientAPICaller interface {
	Start() error
	Init(gateway *UDPAPIGateway)
	run() error
	BodyJson() []byte
}

type CallMessTemp struct {
	FName string `json:"f_name"`
}

type CallerTemp struct {
	Gateway     *UDPAPIGateway
	destIP      []net.UDPAddr
	destID      []string
	CliAPIName  string
	MessContent interface{}
}

func (m *CallerTemp) Init(gateway *UDPAPIGateway) {
	m.Gateway = gateway
}

func (m *CallerTemp) Run() error {
	var ids []net.UDPAddr

	for _, id := range m.destID {
		c, err := clientManage.Get(id)
		if err != nil {
			return err
		}
		ids = append(ids, net.UDPAddr{
			IP:   c.IP,
			Port: c.Port,
			Zone: "",
		})
	}

	err := m.Gateway.SendMess(m.BodyJson(), append(ids, m.destIP...)...)
	return err

}

func (m *CallerTemp) BodyJson() []byte {
	mess, _ := json.Marshal(m.MessContent)
	return mess
}

func (m *CallerTemp) MoreDestByIP(IPs ...net.UDPAddr) {
	m.destIP = append(m.destIP, IPs...)
}

func (m *CallerTemp) MoreDestBySysyncID(IDs ...string) error {
	for _, id := range IDs {
		_, e := clientManage.Get(id)
		if e != nil {
			return e
		}
	}
	m.destID = append(m.destID, IDs...)
	return nil
}
