package ClientAPICallers

import (
	"ConfigServer/APIGateway"
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
	HostIP string,
	HostPort int,
) *UpdateHostName {
	return &UpdateHostName{
		CallerTemp: APIGateway.CallerTemp{
			MessContent: UpdateHostNameMessTemp{
				CallMessTemp: APIGateway.CallMessTemp{
					FName: "update_host_name",
				},
				HostIP:   HostIP,
				HostPort: HostPort,
			},
		},
	}
}
