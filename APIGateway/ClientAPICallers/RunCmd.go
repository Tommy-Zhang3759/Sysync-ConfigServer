package ClientAPICallers

import (
	"ConfigServer/APIGateway"
	"net"
)

type RunCmdMessTemp struct {
	APIGateway.CallMessTemp
	Cmd []byte `json:"command"`
}

type RunCmd struct {
	APIGateway.CallerTemp
	MessContent RunCmdMessTemp
}

func NewRunCmd(
	Dest []net.UDPAddr,
	CliAPIName string,
	cmd []byte,
) *RunCmd {
	return &RunCmd{
		CallerTemp: APIGateway.CallerTemp{
			Dest: Dest,
			MessContent: RunCmdMessTemp{
				CallMessTemp: APIGateway.CallMessTemp{
					FName: CliAPIName,
				},
				Cmd: cmd,
			},
		},
	}
}
