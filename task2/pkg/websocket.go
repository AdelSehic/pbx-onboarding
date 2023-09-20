package amigo

import (
	"fmt"
	"net/http"

	websocket "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}

// Structured way of sending data to clients
type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// Data from Amigo which is important to connected clients
type FetchData struct {
	RegisteredDevices int      `json:"regdev"`
	ActiveDevices     int      `json:"activedev"`
	DeviceList        []Device `json:"devicelist"`
	BridgeCount       int      `json:"bridgecount"`
}

type Device struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// AMIData pulls important data from Amigo structure and prepares it for sending
func (ami *Amigo) AMIData() *FetchData {
	data := &FetchData{
		RegisteredDevices: len(ami.Devices),
		ActiveDevices:     ami.Active,
		DeviceList:        make([]Device, 0, len(ami.Devices)),
		BridgeCount:       ami.Bridges,
	}
	for key, value := range ami.Devices {
		data.DeviceList = append(data.DeviceList, Device{Name: key, Status: value})
	}
	return data
}

// ConnectClient upgrades a client which has request to connect to us to a websocket and sends page initialization information through it
func (ami *Amigo) ConnectClient(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // CORS is unimportant here, accept all requests
	ws, err := upgrader.Upgrade(w, r, nil)                            // function from Gorilla WS library
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := ws.WriteJSON(&Message{Type: "setup", Data: ami.AMIData()}); err != nil {
		fmt.Println(err.Error())
		return
	} // send setup information to connected client and check for errors

	ami.Hub.AddClient(ws) // add the client to our websocket hub
}

// StartWS starts listening to incoming connection requests
func (ami *Amigo) StartWS(addr string) {
	http.HandleFunc("/", ami.ConnectClient)
	go http.ListenAndServe(addr, nil)
	ami.Hub.Start()
}
