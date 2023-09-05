package main

import (
	amigo "ami/pkg"
	"log"
)

func main() {
	ami := amigo.NewManager()
	if err := ami.SetConf("10.1.0.228", "8088", "adel", "123"); err != nil {
		log.Fatal(err.Error())
	}

	if err := ami.Login(); err != nil {
		log.Fatal(err.Error())
	}

	ch := ami.EventListener()

	for event := range ch {
		ami.EventHandler(event)
	}
}
