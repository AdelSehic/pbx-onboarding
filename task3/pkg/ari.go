package ari

import (
	"fmt"
	"time"

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

	ari.AddToCall(call, dev...)

	call.Ring()

	time.Sleep(5 * time.Second)

	ari.AddToCall(call, "102")
	call.Ring()

	time.Sleep(5 * time.Second)

	// for chanCount >= min {
	// 	for ch := range chans {
	// 		data, _ := ari.Client.Channel().Data(ch.Key())
	// 		if data == nil {
	// 			delete(chans, ch)
	// 			chanCount--
	// 		}
	// 	}
	// }
}

// func (ari *Ari) directCall(ext1, ext2 string) {
// 	handle1, err := ari.Client.Channel().Create(ari.AppKey, ariLib.ChannelCreateRequest{
// 		Endpoint: exten[ext1],
// 		App:      ari.AppName,
// 	})
// 	if err != nil {
// 		fmt.Printf("Error creating a channel to %s endpoint\n", ext1)
// 		return
// 	}
// 	handle2, err := ari.Client.Channel().Create(ari.AppKey, ariLib.ChannelCreateRequest{
// 		Endpoint: exten[ext1],
// 		App:      ari.AppName,
// 	})
// 	if err != nil {
// 		fmt.Printf("Error creating a channel to %s endpoint\n", ext2)
// 		return
// 	}
// }
