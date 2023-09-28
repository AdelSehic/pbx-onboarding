package ari

import (
	"fmt"
	"log"

	ariLib "github.com/CyCoreSystems/ari"
	ariClient "github.com/CyCoreSystems/ari/client/native"
)

type Ari struct {
	AppName string
	AppKey  *ariLib.Key
	Client  ariLib.Client
	Calls   map[string]*Call
}

type Call struct {
	ID         string
	Conference bool
	Channels   map[*ariLib.ChannelHandle]struct{}
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

	min := 0
	if len(dev) == 2 {
		min = 2
	} else if len(dev) > 2 {
		min = 1
	} else {
		fmt.Println("must enter at least two exensions to dial")
		return
	}

	chanCount := 0
	chans := make(map[*ariLib.ChannelHandle]struct{}, len(dev))
	bridge, err := ari.Client.Bridge().Create(ari.AppKey, "", "")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer func() { // exit cleanup
		for ch := range chans {
			ch.Hangup()
		}
		bridge.Delete()
	}()

	for i := range dev {
		handle, err := ari.Client.Channel().Create(ari.AppKey, ariLib.ChannelCreateRequest{
			Endpoint: exten[dev[i]],
			App:      ari.AppName,
		})
		if err != nil {
			fmt.Printf("Error creating a channel to %s endpoint\n", dev[i])
			if min == 2 {
				return
			}
			continue
		}
		if err := handle.Dial("", 15); err != nil {
			fmt.Printf("Error dialing %s\n", dev[i])
			if min == 2 {
				return
			}
			continue
		}
		bridge.AddChannel(handle.ID())
		chans[handle] = struct{}{}
		chanCount++
	}

	for chanCount >= min {
		for ch := range chans {
			data, _ := ari.Client.Channel().Data(ch.Key())
			if data == nil {
				delete(chans, ch)
				chanCount--
			}
		}
	}
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
