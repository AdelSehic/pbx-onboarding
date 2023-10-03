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
	ToRing     chan *ariLib.ChannelHandle
	MinActive  int
}

func (ari *Ari) NewCall() (*Call, error) {

	bridge, err := ari.Client.Bridge().Create(ari.AppKey.New("", ""), "", "")
	if err != nil {
		return nil, err
	}

	call := &Call{
		ID:         bridge.ID(),
		Bridge:     bridge,
		Conference: false,
		ChanCount:  0,
		Channels:   make(map[*ariLib.ChannelHandle]struct{}),
		ToRing:     make(chan *ariLib.ChannelHandle, 16),
		MinActive:  2,
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
		call.ToRing <- handle
		call.Bridge.AddChannel(handle.ID())
		call.Channels[handle] = struct{}{}
		call.ChanCount++
		if call.ChanCount > 2 {
			call.Conference = true
			call.MinActive = 1
		}
	}
}

func (call *Call) Ring() {
	for ch := range call.ToRing {
		if err := ch.Dial("Asterisk REST interface", 15); err != nil {
			fmt.Printf("error on ringing %s\n", ch.Key().ID)
			continue
		}
	}
}

func (ari *Ari) MonitorCall(call *Call) {
	for call.ChanCount >= call.MinActive {
		for ch := range call.Channels {
			data, _ := ari.Client.Channel().Data(ch.Key())
			if data == nil {
				delete(call.Channels, ch)
				call.ChanCount--
			}
		}
	}
}

func (ari *Ari) CloseCall(call *Call) {
	for ch := range call.Channels {
		ch.Hangup()
	}
	call.Bridge.Delete()
	delete(ari.Calls, call.ID)
}
