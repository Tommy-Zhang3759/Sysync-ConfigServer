package WebAPI

type WebSocketMessagePusher interface {
	Start() error
	Stop() error
	NewMessage() error

	pushMessage() error
}

type WebSocketMessagePusherImpl struct {
	messageQueue chan string
}

func (pusher *WebSocketMessagePusherImpl) f() {

}
