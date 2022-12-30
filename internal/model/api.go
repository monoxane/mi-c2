package model

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebsocketConnection struct {
	Mux     *sync.Mutex
	Socket  *websocket.Conn
	Cluster string
}
