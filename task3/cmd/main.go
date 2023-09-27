package main

import (
	"fmt"
	"log"

	ari "github.com/CyCoreSystems/ari"
	ariClient "github.com/CyCoreSystems/ari/client/native"
)

func main() {
	client, err := ariClient.Connect(&ariClient.Options{
		Application:  "blondie",
		URL:          "http://10.1.0.228:8088/ari",
		WebsocketURL: "ws://10.1.0.228:8088/ari/events",
		Username:     "asterisk",
		Password:     "test123",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	key := ari.AppKey(client.ApplicationName())

	keys, err := client.Endpoint().List(key)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, key := range keys {
		data, _ := client.Endpoint().Data(key)
		fmt.Println(data)
	}

	// client.Bridge().Create()
}
