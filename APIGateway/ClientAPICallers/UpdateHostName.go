package ClientAPICallers

import (
	"ConfigServer/APIGateway"
	"net"
)

type UpdateHostNameMessTemp struct {
	APIGateway.CallMessTemp
	HostIP   string `json:"host_ip,omitempty"`
	HostPort int    `json:"host_port,omitempty"`
}

type UpdateHostName struct {
	APIGateway.CallerTemp
	MessContent UpdateHostNameMessTemp
}

func NewUpdateHostName(
	Dest []net.UDPAddr,
	CliAPIName string,
	HostIP string,
	HostPort int,
) *UpdateHostName {
	return &UpdateHostName{
		CallerTemp: APIGateway.CallerTemp{
			Dest: Dest,
			MessContent: UpdateHostNameMessTemp{
				CallMessTemp: APIGateway.CallMessTemp{
					FName: CliAPIName,
				},
				HostIP:   HostIP,
				HostPort: HostPort,
			},
		},
	}
}
