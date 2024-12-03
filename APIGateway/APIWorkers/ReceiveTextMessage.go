package APIWorkers

import (
	"ConfigServer/APIGateway"
	"ConfigServer/clientManage"
	"encoding/json"
	"fmt"
)

type MessFromClientReceiver struct {
	APIGateway.UDPAPIPortTemp
	cliContainer *clientManage.CliContainer
	messTube     chan ClientMessage // run as a queue to return message to the interface connected to the front end
}

type ClientMessage struct {
	Error    string `json:"error,omitempty"`
	Message  string `json:"message,omitempty"`
	SysyncID string `json:"sysyncID"`
}

func NewClientMessageError(err string, sysyncID string) ClientMessage {
	return ClientMessage{
		Error:    err,
		SysyncID: sysyncID,
	}
}

func NewClientMessage(message string, sysyncID string) ClientMessage {
	return ClientMessage{
		Message:  message,
		SysyncID: sysyncID,
	}
}

func (message *ClientMessage) String() (string, error) {
	if message.Error == "" && message.Message == "" {
		return "", fmt.Errorf("empty message")
	} else if message.Error != "" && message.Message != "" {
		return "", fmt.Errorf("invalid message: %s", message.Message)
	} else {
		messJson, _ := json.Marshal(message)
		return string(messJson), nil
	}
}

func NewMessFromClientReceiver(cliContainer *clientManage.CliContainer, bufferLen ...int) *MessFromClientReceiver {
	if len(bufferLen) == 0 {
		bufferLen = append(bufferLen, 1)
	}
	return &MessFromClientReceiver{
		cliContainer: cliContainer,
		messTube:     make(chan ClientMessage, bufferLen[0]),
	}
}

func (u *MessFromClientReceiver) ReadMessage() string {
	var err = fmt.Errorf("e")
	var message string
	for err != nil {
		clientMessage := <-u.messTube
		message, err = clientMessage.String()
	}
	return message
}

func (u *MessFromClientReceiver) Start() error {
	stop := false

	for stop == false {
		reqPack := u.MessageQue.Pop().(APIGateway.UDPMessage)

		select {
		case <-u.EndRun:
			fmt.Println("Received stop signal, goroutine exiting...")
			stop = true
		default:
			type ClientRequest struct {
				TaskID   string `json:"taskID"`
				SysyncID string `json:"sysync_id"`
				Error    string `json:"error,omitempty"`
				Message  string `json:"message,omitempty"`
			}

			go func() {
				rsp := APIGateway.ApiResponse{
					Fname:   u.GetKeyWord(),
					Status:  100,
					Message: "",
					Error:   "",
				}

				rsp.Fname = u.GetKeyWord()

				var keyExist bool
				var sysyncID, cliErr, cliMess string
				if sysyncID, keyExist = reqPack.Text["sysync_id"].(string); keyExist == false {
					rsp.Error = "sysync_id dose not exist"
					rsp.Status = 400
					rspJson, _ := json.Marshal(rsp)
					_ = u.Gateway.SendMess(rspJson, reqPack.Addr)
					return
				}

				if cliErr, keyExist = reqPack.Text["error"].(string); keyExist == true { // as an error message
					u.messTube <- NewClientMessageError(cliErr, sysyncID)
				} else if cliMess, keyExist = reqPack.Text["message"].(string); keyExist == true { // as an info message
					u.messTube <- NewClientMessage(cliMess, sysyncID)
				} else {
					rsp.Error = "invalid message"
					rsp.Status = 400
					rspJson, _ := json.Marshal(rsp)
					_ = u.Gateway.SendMess(rspJson, reqPack.Addr)
					return
				}

				rsp.Status = 200
				rsp.Message = "success"
				rspJson, _ := json.Marshal(rsp)
				_ = u.Gateway.SendMess(rspJson, reqPack.Addr)
			}()
		}
	}
	return nil
}
