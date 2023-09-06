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
	if err := ami.SetConf("10.1.0.228", "8088", "adel", "123"); err != nil {
		log.Fatal(err.Error())
	}

	if err := ami.Login(); err != nil {
		log.Fatal(err.Error())
	}
	ch := ami.EventListener()

	// ctx, cancel := context.WithCancel(context.Background())

loop:
	for {
		select {
		case event := <-ch:
			ami.EventHandler(event)
		case <-stop:
			ami.ListDevices()
			break loop
		}
	}
}
