package amigo

import (
	websocket "github.com/gorilla/websocket"
)

// WSHub acts as a very simple communcation interface to enable us multicasting in case multiple clients want to listen to us
type WSHub struct {
	Clients   map[*websocket.Conn]struct{}
	Broadcast chan Message
}

// NewHub initializes WSHub structure
func NewHub() *WSHub {
	return &WSHub{
		Clients:   make(map[*websocket.Conn]struct{}),
		Broadcast: make(chan Message),
	}
}

// AddClients adds a new client to our multicast
func (ws *WSHub) AddClient(cl *websocket.Conn) {
	ws.Clients[cl] = struct{}{}
}

// Start starts a goroutine which waits for an incoming message and sends it to all connected clients, effectively creating multicast
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
