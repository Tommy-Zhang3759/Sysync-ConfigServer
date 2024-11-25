package ClientAPICallers

import (
	"ConfigServer/APIGateway"
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
	cmd []byte,
) *RunCmd {
	return &RunCmd{
		CallerTemp: APIGateway.CallerTemp{
			MessContent: RunCmdMessTemp{
				CallMessTemp: APIGateway.CallMessTemp{
					FName: "run_command",
				},
				Cmd: cmd,
			},
		},
	}
}
