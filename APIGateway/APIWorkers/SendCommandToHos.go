package APIWorkers

import "ConfigServer/APIGateway"

type SendCommandToHost struct {
	APIGateway.UDPAPIPortTemp
}

func (h *SendCommandToHost) Run() error {
	return nil
}
