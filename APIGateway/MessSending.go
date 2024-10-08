package APIGateway

import (
	"encoding/json"
	"net"
)

type MessSending struct {
	UDPAPIPortTemp
	Dest        []net.UDPAddr
	MessContent map[string]interface{}
}

func (m *MessSending) name() {

}
func (m *MessSending) Run() error {
	err := m.Gateway.SendMess(m.bodyJson(), m.Dest...)
	return err
}

func (m *MessSending) bodyJson() []byte {
	mess, _ := json.Marshal(m.MessContent)
	return mess
}
