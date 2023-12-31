package ari

import (
	"context"
	"fmt"
	"sync"

	ariLib "github.com/CyCoreSystems/ari"
	ariClient "github.com/CyCoreSystems/ari/client/native"
)

// Ari struct for holding all the necessarhy information for call monitoring and managing
type Ari struct {
	AppName string
	AppKey  *ariLib.Key
	Client  ariLib.Client
	Calls   map[string]*Call
	Wg      *sync.WaitGroup
}

// New initializes an Ari struct, requires a configuration to connect to Asterisk right away
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
		Wg:      &sync.WaitGroup{},
	}, nil
}

// Dial creates a new call with specified clients added right away
func (ari *Ari) Dial(ctx context.Context, dev ...string) {

	call, err := ari.NewCall()
	if err != nil {
		fmt.Println("Failed to create a new call")
		return
	}
	ari.Calls[call.ID] = call

	ari.AddToCall(call, dev...)
	go ari.MonitorCall(ctx, call)
}

// List prints all active calls and their participants
func (ari *Ari) List() {
	for _, c := range ari.Calls {
		fmt.Println(c.ID)
		for _, devs := range c.Channels {
			data, _ := devs.Data()
			fmt.Printf("\t- %s : %s\n", data.Name, data.Key)
		}
	}
}
