package main

import (
	amigo "ami/pkg"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ami := amigo.NewManager()
	ami.SetConf("192.168.0.11", "5038", "adel", "123")
	evChan, err := ami.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	// evChan, errChan := ami.EventListener()
	// ami.StartWS()

loop:
	for {
		select {
		case event := <-evChan:
			ami.EventHandler(event)
		case <-stop:
			break loop
		}
	}
}
