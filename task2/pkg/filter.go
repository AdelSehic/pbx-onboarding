package amigo

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// trimInfo is used to pull data from responses such as "Event: <event name>". It would take event name and return it to us.
func trimInfo(line string) string {
	return strings.Trim(strings.Split(line, " ")[1], "\r\n")
}

// a simple way to check if we have a function defined for an event and call it if we do
var filter = map[string]func([]string, *Amigo){
	"SuccessfulAuth":    succAuth,
	"DeviceStateChange": devStateChange,
	"BridgeCreate":      addBridge,
	"BridgeDestroy":     rmBridge,
}

// rest of this file are event handling functions

func succAuth(event []string, ami *Amigo) {
	dev, ip := trimInfo(event[6]), trimInfo(event[9])
	log.Printf("Successful Authentication by %s from %s\n", dev, ip)

	ami.Hub.Broadcast <- Message{ // send data to websockets connected to us
		Type: "succauth",
		Data: []string{dev, ip},
	}
}

func devStateChange(event []string, ami *Amigo) {
	dev, state := trimInfo(event[2]), trimInfo(event[3])
	prevstate := ami.Devices[dev]
	ami.Devices[dev] = state
	log.Printf("%s is now %s", dev, state)

	if prevstate == "UNAVAILABLE" {
		ami.Active++
	}
	if state == "UNAVAILABLE" {
		ami.Active--
	} // previous state hos to be checked to determine the amount of active devices

	ami.Hub.Broadcast <- Message{
		Type: "devstatechange",
		Data: []string{dev, ami.Devices[dev]},
	}

	ami.Hub.Broadcast <- Message{
		Type: "activedevs",
		Data: ami.Active,
	}
}

func addBridge(event []string, ami *Amigo) {
	ami.Bridges++
	fmt.Println("Bridge created")

	ami.Hub.Broadcast <- Message{
		Type: "brcountupdate",
		Data: []string{"Bridge created", strconv.Itoa(ami.Bridges)},
	}
}

func rmBridge(event []string, ami *Amigo) {
	ami.Bridges--
	fmt.Println("Bridge destroyed")

	ami.Hub.Broadcast <- Message{
		Type: "brcountupdate",
		Data: []string{"Bridge destroyed", strconv.Itoa(ami.Bridges)},
	}
}
