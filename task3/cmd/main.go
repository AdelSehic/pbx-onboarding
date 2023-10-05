package main

import (
	ari "ari/pkg"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	ari, err := ari.New("blondie", "10.1.0.228:8088", "asterisk", "test123")
	if err != nil {
		log.Fatal(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		args := getInput(*reader)

		switch args[0] {
		case "dial":
			ari.Dial(args[1:]...)
		case "list":
			ari.List()
		case "join":
			ari.JoinCall(args)
		default:
			fmt.Println("Invalid option")
		}
	}
}

func getInput(reader bufio.Reader) []string {
	fmt.Scanln(reader)
	read, _ := reader.ReadString('\n')
	input := strings.Trim(read, "\n")
	args := strings.Split(input, " ")
	if len(args) < 1 {
		return nil
	}
	return args
}
