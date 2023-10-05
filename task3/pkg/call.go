package ari

import (
	"fmt"

	ariLib "github.com/CyCoreSystems/ari"
)

type Call struct {
	ID        string
	ChanCount int
	Bridge    *ariLib.BridgeHandle
	Channels  map[*ariLib.ChannelHandle]struct{}
	ToRing    chan *ariLib.ChannelHandle
	MinActive int
}

const maxParticipants = 16
const globalTimeout = 15

func (ari *Ari) NewCall() (*Call, error) {

	bridge, err := ari.Client.Bridge().Create(ari.AppKey.New("", ""), "", "")
	if err != nil {
		return nil, err
	}

	call := &Call{
		ID:        bridge.ID(),
		Bridge:    bridge,
		ChanCount: 0,
		Channels:  make(map[*ariLib.ChannelHandle]struct{}),
		ToRing:    make(chan *ariLib.ChannelHandle, maxParticipants),
		MinActive: 2,
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
			call.MinActive = 1
		}
	}
}

func (call *Call) Ring() {
	for ch := range call.ToRing {
		if err := ch.Dial("Asterisk REST interface", globalTimeout); err != nil {
			fmt.Printf("error on ringing %s\n", ch.Key().ID)
			continue
		}
	}
}

func (ari *Ari) JoinCall(args []string) {
	if len(args) <= 2 {
		fmt.Println(`invalid format, propper format is "join <callid> clients..." `)
		return
	}
	if _, ok := ari.Calls[args[1]]; !ok {
		fmt.Println("specified call ID does not exist")
		return
	}
	ari.AddToCall(ari.Calls[args[1]], args[2:]...)
}

func (ari *Ari) MonitorCall(call *Call) {

	sub := call.Bridge.Subscribe(ariLib.Events.ChannelLeftBridge).Events()

	for {
		<-sub
		call.ChanCount--
		if call.ChanCount < call.MinActive {
			break
		}
	}

	ari.CloseCall(call)
}

func (ari *Ari) CloseCall(call *Call) {
	for ch := range call.Channels {
		ch.Hangup()
	}
	call.Bridge.Delete()
	delete(ari.Calls, call.ID)
}
