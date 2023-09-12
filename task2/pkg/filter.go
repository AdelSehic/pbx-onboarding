package amigo

import (
	"fmt"
	"log"
	"strings"
)

func trimInfo(line string) string {
	return strings.Trim(strings.Split(line, " ")[1], "\r\n")
}

var filter = map[string]func([]string, *Amigo){
	"SuccessfulAuth":    succAuth,
	"DeviceStateChange": devStateChange,
	"BridgeCreate":      addBridge,
	"BridgeDestroy":     rmBridge,
}

func succAuth(event []string, ami *Amigo) {
	dev, ip := trimInfo(event[6]), trimInfo(event[9])
	log.Printf("Successful Authentication by %s from %s\n", dev, ip)
}

func devStateChange(event []string, ami *Amigo) {
	dev, state := trimInfo(event[2]), trimInfo(event[3])
	log.Printf("%s is now %s\n", dev, state)

	if trimInfo(event[3]) == "NOT_INUSE" {
		ami.Active++
	} else if ami.Active > 0 {
		ami.Active--
	}

	ami.Hub.Broadcast <- Message{
		Type: "devstatechange",
		Data: ami.Devices[dev],
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
		Data: ami.Bridges,
	}
}

func rmBridge(event []string, ami *Amigo) {
	ami.Bridges--
	fmt.Println("Bridge destroyed")
	ami.Hub.Broadcast <- Message{
		Type: "brcountupdate",
		Data: ami.Bridges,
	}
}
