package ari

import (
	"context"
	"fmt"
	"sync"

	"github.com/CyCoreSystems/ari"
	ariLib "github.com/CyCoreSystems/ari"
)

const globalTimeout = 15
const staticArguments = 2
const maxClientsToCall = 10

// Call struct that holds necessary information for calls to be managed to propperly
type Call struct {
	ID                        string
	ChanCount                 int
	Bridge                    *ariLib.BridgeHandle
	Channels                  map[string]*ariLib.ChannelHandle
	MinimumActiveParticipants int
	mu                        sync.Mutex
}

// NewCall initializes a call struct
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
		mu:                        sync.Mutex{},
	}

	return call, nil
}

// AddToCall takes an existing call and adds new clients to it
func (ari *Ari) AddToCall(call *Call, dev ...string) {

	devs := make([]*ariLib.ChannelHandle, 0, maxClientsToCall)

	call.mu.Lock()
	defer call.mu.Unlock()

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

// Ring rings all channels that were recently added but not called yet
func (call *Call) Ring(devs []*ari.ChannelHandle) {
	for _, ch := range devs {
		if err := ch.Dial("Asterisk REST interface", globalTimeout); err != nil {
			fmt.Printf("error on ringing %s\n", ch.Key().ID)
			continue
		}
	}
}

// JoinCall takes existing call's ID as a first parameter and adds all specified clients to it
func (ari *Ari) JoinCall(args []string) {
	if len(args) <= staticArguments {
		fmt.Println(`invalid format, propper format is "join <callid> clients..." `)
		return
	}
	if _, ok := ari.Calls[args[1]]; !ok {
		fmt.Println("specified call ID does not exist")
		return
	}
	ari.AddToCall(ari.Calls[args[1]], args[2:]...)
}

// MonitorCall - the main function of our program, it monitors a call for leave events and closes it when it makes sense
func (ari *Ari) MonitorCall(ctx context.Context, call *Call) {

	ari.Wg.Add(1) // add another call monitor to waitgroup sync
	sub := call.Bridge.Subscribe(ariLib.Events.ChannelLeftBridge).Events()

loop:
	for {
		select {
		case event := <-sub:
			call.mu.Lock()
			call.ChanCount--
			delete(call.Channels, event.Keys().First().ID)
			call.mu.Unlock()
			if call.ChanCount < call.MinimumActiveParticipants {
				break loop
			}
		case <-ctx.Done():
			fmt.Printf("Breaking the call %s\n", call.ID)
			break loop
		}
	}

	ari.CloseCall(call)
}

// CloseCall shuts down the specified call cleanly
func (ari *Ari) CloseCall(call *Call) {

	call.mu.Lock()
	defer call.mu.Unlock()

	for _, ch := range call.Channels {
		ch.Hangup()
	}
	call.Bridge.Delete()
	delete(ari.Calls, call.ID)
	ari.Wg.Done()
}
