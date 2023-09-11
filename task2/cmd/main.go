package main

import (
	amigo "ami/pkg"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ami := amigo.NewManager()
	if err := ami.SetConf("192.168.0.11", "8088", "adel", "123"); err != nil {
		log.Fatal(err.Error())
	}

	if err := ami.Login(); err != nil {
		log.Fatal(err.Error())
	}
	evChan, errChan := ami.EventListener()

	http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) {
		data, err := ami.MarshallAMI()
		if err != nil {
			http.Error(w, "Failed to write JSON", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})

	go http.ListenAndServe(":9999", nil)

loop:
	for {
		select {
		case event := <-evChan:
			ami.EventHandler(event)
		case err := <-errChan:
			log.Fatal(err.Error())
		case <-stop:
			ami.ListDevices()
			fmt.Printf("Number of active calls: %d\n", ami.Bridges)
			break loop
		}
	}
}
