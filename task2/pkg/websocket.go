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

type Message struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

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

func (ami *Amigo) ConnectClient(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	message := &Message{
		Type: "setup",
		Data: ami.AMIData(),
	}
	if err := ws.WriteJSON(message); err != nil {
		fmt.Println(err.Error())
		return
	}
	ami.Hub.AddClient(ws)
}

func (ami *Amigo) StartWS() {
	http.HandleFunc("/", ami.ConnectClient)
	go http.ListenAndServe(":9999", nil)
	ami.Hub.Start()
}
