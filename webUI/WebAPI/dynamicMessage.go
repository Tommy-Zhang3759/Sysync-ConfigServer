package WebAPI

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var webSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := webSocketUpgrader.Upgrade(w, r, nil)
	return conn, err
}
