package amigo

import (
	"encoding/json"
)

type FetchData struct {
	DeviceCount int
	BridgeCount int
	DeviceList  []Device
}

type Device struct {
	Name   string
	Status string
}

func (ami *Amigo) MarshallAMI() ([]byte, error) {
	data := &FetchData{
		BridgeCount: ami.Bridges,
		DeviceCount: len(ami.Devices),
		DeviceList:  make([]Device, 0, len(ami.Devices)),
	}
	for key, value := range ami.Devices {
		data.DeviceList = append(data.DeviceList, Device{Name: key, Status: value})
	}
	post, err := json.Marshal(data)
	return post, err
}
