package APIGateway

import (
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
	Dest        []net.UDPAddr
	CliAPIName  string
	MessContent interface{}
}

func (m *CallerTemp) Start() error {
	go func() {
		_ = m.run()
	}()
	return nil
}

func (m *CallerTemp) Init(gateway *UDPAPIGateway) {
	m.Gateway = gateway
}

func (m *CallerTemp) run() error {
	err := m.Gateway.SendMess(m.BodyJson(), m.Dest...)
	return err

}

func (m *CallerTemp) BodyJson() []byte {
	mess, _ := json.Marshal(m.MessContent)
	return mess
}
