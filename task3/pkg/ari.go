package ari

import (
	"fmt"

	ariLib "github.com/CyCoreSystems/ari"
	ariClient "github.com/CyCoreSystems/ari/client/native"
)

type Ari struct {
	AppName string
	AppKey  *ariLib.Key
	Client  ariLib.Client
	Calls   map[string]*Call
}

func New(appName, address, user, password string) (*Ari, error) {
	client, err := ariClient.Connect(&ariClient.Options{
		Application:  appName,
		URL:          "http://" + address + "/ari",
		WebsocketURL: "ws://" + address + "/ari/events",
		Username:     user,
		Password:     password,
	})
	if err != nil {
		return nil, err
	}
	appKey := ariLib.AppKey(client.ApplicationName())

	return &Ari{
		AppName: appName,
		AppKey:  appKey,
		Client:  client,
		Calls:   make(map[string]*Call),
	}, nil
}

func (ari *Ari) Dial(dev ...string) {

	call, err := ari.NewCall()
	if err != nil {
		fmt.Println("Failed to create a new call")
		return
	}
	defer call.Close()
	ari.Calls[call.ID] = call
	fmt.Println("call id: ", call.ID)

	ari.AddToCall(call, dev...)
	go call.Ring()

	ari.MonitorCall(call)
	call.Close()
}

func (ari *Ari) List() {
	for _, c := range ari.Calls {
		fmt.Println(c.ID)
		for devs := range c.Channels {
			data, _ := devs.Data()
			fmt.Printf("\t%s\n", data.Name)
		}
	}
}
