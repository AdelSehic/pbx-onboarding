package ari

import (
	"fmt"

	ariLib "github.com/CyCoreSystems/ari"
)

type Call struct {
	ID         string
	Conference bool
	ChanCount  int
	Bridge     *ariLib.BridgeHandle
	Channels   map[*ariLib.ChannelHandle]struct{}
}

func (ari *Ari) NewCall() (*Call, error) {

	bridge, err := ari.Client.Bridge().Create(ari.AppKey, "", "")
	if err != nil {
		return nil, err
	}

	call := &Call{
		ID:         bridge.Key().ID,
		Bridge:     bridge,
		Conference: false,
		ChanCount:  0,
		Channels:   make(map[*ariLib.ChannelHandle]struct{}),
	}

	return call, nil
}

func (ari *Ari) AddToCall(call *Call, dev ...string) {
	for i := range dev {
		handle, err := ari.Client.Channel().Create(ari.AppKey, ariLib.ChannelCreateRequest{
			Endpoint: exten[dev[i]],
			App:      ari.AppName,
		})
		if err != nil {
			fmt.Printf("Error creating a channel to %s endpoint\n", dev[i])
			continue
		}
		call.Bridge.AddChannel(handle.ID())
		call.Channels[handle] = struct{}{}
		call.ChanCount++
	}
}

func (call *Call) Ring() {
	for ch := range call.Channels {
		if err := ch.Dial("Asterisk REST interface", 15); err != nil {
			fmt.Printf("error on ringing %s\n", ch.Key().ID)
		}
	}
}

func (call *Call) Close() {
	for ch := range call.Channels {
		ch.Hangup()
	}
	call.Bridge.Delete()
}
