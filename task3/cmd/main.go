package main

import (
	config "ari/configs"
	ari "ari/pkg"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	interrupt := make(chan os.Signal, 1) // hijacks kill singals so we can break our program cleanly
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	ctx, stop := context.WithCancel(context.Background()) // context for propagating stop signal to functions

	cfg, err := config.LoadConfig("../configs/config.json") // load config from file
	if err != nil {
		log.Fatal(err.Error())
	}

	ari, err := ari.New(cfg.GetConfig()) // send loaded config to ari
	if err != nil {
		log.Fatal(err.Error())
	}

	input := getInput() // define stdin variable
loop:
	for {
		select {
		case args := <-input:
			switch args[0] {
			case "dial":
				ari.Dial(ctx, args[1:]...)
			case "list":
				ari.List()
			case "join":
				ari.JoinCall(args)
			default:
				fmt.Println("Invalid option")
			}
		case <-interrupt:
			fmt.Println("Closing the application....")
			stop()
			break loop
		}
	}
	ari.Wg.Wait() // wait for all calls to be stopped before closing out
}

func getInput() chan []string {

	ch := make(chan []string)
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			fmt.Scanln(reader)
			read, _ := reader.ReadString('\n')
			input := strings.Trim(read, "\n")
			args := strings.Split(input, " ")
			if len(args) < 1 {
				continue
			}
			ch <- args
		}
	}()
	return ch
}
