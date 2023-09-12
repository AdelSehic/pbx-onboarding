package amigo

import (
	"log"
	"regexp"
)

var filter = map[string]func(string, *Amigo){
	"SuccessfulAuth":    succAuth,
	"DeviceStateChange": devStateChange,
	"BridgeCreate":      addBridge,
	"BridgeDestroy":     rmBridge,
}

func succAuth(event string, ami *Amigo) {
	regID := regexp.MustCompile(`AccountID: (\w+)`)
	id := regID.FindStringSubmatch(event)[1]

	regIP := regexp.MustCompile(`RemoteAddress: \w+/\w+/([\d\.]+)/.+`)
	ip := regIP.FindStringSubmatch(event)[1]

	log.Printf("Successful Authentication by %s from %s\n", id, ip)
}

func devStateChange(event string, ami *Amigo) {
	reg := regexp.MustCompile(`: ([^\r\n]+)`)
	caught := reg.FindAllStringSubmatch(event, -1)

	log.Printf("%s is now %s\n", caught[4][1], caught[5][1])
	ami.Devices[caught[4][1]] = caught[5][1]

	ami.Hub.Broadcast <- Message{
		Type: "devstatechange",
		Data: ami.Devices[caught[4][1]],
	}
	// ami.FetchBridgeCount()
}

func addBridge(event string, ami *Amigo) {
	ami.Bridges++
	ami.Hub.Broadcast <- Message{
		Type: "brcountupdate",
		Data: ami.Bridges,
	}
}

func rmBridge(event string, ami *Amigo) {
	ami.Bridges--
	ami.Hub.Broadcast <- Message{
		Type: "brcountupdate",
		Data: ami.Bridges,
	}
}
