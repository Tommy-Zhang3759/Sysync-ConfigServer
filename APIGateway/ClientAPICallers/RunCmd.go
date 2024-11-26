package ClientAPICallers

import (
	"ConfigServer/APIGateway"
)

type RunCmdMessTemp struct {
	APIGateway.CallMessTemp
	Cmd string `json:"command"`
}

type RunCmd struct {
	APIGateway.CallerTemp
	MessContent RunCmdMessTemp
}

func NewRunCmd(
	cmd string,
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
