package APIWorkers

import (
	"ConfigServer/APIGateway"
	"encoding/json"
	"net"
)

type TextMessage struct {
	APIGateway.UDPAPIPortTemp
	Dest        []net.UDPAddr
	MessContent interface{}
}

func (m *TextMessage) name() {
	return
}

func (m *TextMessage) Run() error {
	err := m.Gateway.SendMess(m.bodyJson(), m.Dest...)
	return err
}

func (m *TextMessage) bodyJson() []byte {
	mess, _ := json.Marshal(m.MessContent)
	return mess
}
