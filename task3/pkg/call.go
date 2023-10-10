package ari

import (
	"fmt"

	"github.com/CyCoreSystems/ari"
	ariLib "github.com/CyCoreSystems/ari"
)

type Call struct {
	ID                        string
	ChanCount                 int
	Bridge                    *ariLib.BridgeHandle
	Channels                  map[string]*ariLib.ChannelHandle
	MinimumActiveParticipants int
}

const globalTimeout = 15

func (ari *Ari) NewCall() (*Call, error) {

	bridge, err := ari.Client.Bridge().Create(ari.AppKey.New("", ""), "", "")
	if err != nil {
		return nil, err
	}

	call := &Call{
		ID:                        bridge.ID(),
		Bridge:                    bridge,
		ChanCount:                 0,
		Channels:                  make(map[string]*ariLib.ChannelHandle),
		MinimumActiveParticipants: 2,
	}

	return call, nil
}

func (ari *Ari) AddToCall(call *Call, dev ...string) {

	devs := make([]*ariLib.ChannelHandle, 0, 10)

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
		call.Channels[handle.ID()] = handle
		devs = append(devs, handle)
		call.ChanCount++
		if call.ChanCount > 2 {
			call.MinimumActiveParticipants = 1
		}
	}

	if call.ChanCount < call.MinimumActiveParticipants {
		fmt.Println("Not enough participants to start the call, aborting")
		return
	}
	call.Ring(devs)
}

func (call *Call) Ring(devs []*ari.ChannelHandle) {
	for _, ch := range devs {
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
		event := <-sub
		call.ChanCount--
		delete(call.Channels, event.Keys().First().ID)
		if call.ChanCount < call.MinimumActiveParticipants {
			break
		}
	}

	ari.CloseCall(call)
}

func (ari *Ari) CloseCall(call *Call) {
	for _, ch := range call.Channels {
		ch.Hangup()
	}
	call.Bridge.Delete()
	delete(ari.Calls, call.ID)
}
