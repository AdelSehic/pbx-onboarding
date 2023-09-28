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
	}, nil
}

func (ari *Ari) Call(dev ...string) {

	min := 0
	if len(dev) == 2 {
		min = 2
	} else if len(dev) > 2 {
		min = 1
	} else {
		log.Fatal("must enter at least two numbers to call")
	}

	chanCount := 0
	chans := make(map[*ariLib.ChannelHandle]struct{}, len(dev))
	bridge, err := ari.Client.Bridge().Create(ari.AppKey, "", "conf1")
	if err != nil {
		log.Fatal(err.Error())
	}

	for i := range dev {
		handle, err := ari.Client.Channel().Create(ari.AppKey, ariLib.ChannelCreateRequest{
			Endpoint: dev[i],
			App:      ari.AppName,
		})
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if err := handle.Dial("", 15); err != nil {
			fmt.Println(err.Error())
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

	for ch := range chans {
		ch.Hangup()
	}
	bridge.Delete()
}
