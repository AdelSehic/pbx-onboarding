package main

import (
	amigo "ami/pkg"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	stop := make(chan os.Signal, 1) // hijacks kill singals so we can break our program cleanly
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ami := amigo.NewManager()
	ami.SetConf("10.1.0.228", "5038", "adel", "123")

	evChan, err := ami.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	ami.StartWS(":9999")

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
