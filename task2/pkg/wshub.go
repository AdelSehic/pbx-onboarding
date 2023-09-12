package amigo

import (
	websocket "github.com/gorilla/websocket"
)

type WSHub struct {
	Clients   map[*websocket.Conn]struct{}
	Broadcast chan Message
}

func NewHub() *WSHub {
	return &WSHub{
		Clients:   make(map[*websocket.Conn]struct{}),
		Broadcast: make(chan Message),
	}
}

func (ws *WSHub) AddClient(cl *websocket.Conn) {
	ws.Clients[cl] = struct{}{}
}

func (ws *WSHub) Start() {
	go func() {
		for {
			msg := <-ws.Broadcast
			for client := range ws.Clients {
				client.WriteJSON(msg)
			}
		}
	}()
}
