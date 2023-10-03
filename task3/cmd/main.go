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

	conn, err := ari.New("blondie", "10.1.0.228:8088", "asterisk", "test123")
	if err != nil {
		log.Fatal(err.Error())
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Scanln(reader)
		read, _ := reader.ReadString('\n')
		input := strings.Trim(read, "\n")
		args := strings.Split(input, " ")
		if len(args) < 1 {
			fmt.Println("bad input")
			continue
		}
		switch args[0] {
		case "dial":
			go conn.Dial(args[1:]...)
		case "list":
			conn.List()
		case "join":
			if len(args) <= 2 {
				fmt.Println(`invalid format, propper format is "join <callid> clients..." `)
				continue
			}
			if _, ok := conn.Calls[args[1]]; !ok {
				fmt.Println("specified call ID does not exist")
				continue
			}
			conn.AddToCall(conn.Calls[args[1]], args[2:]...)
		default:
			fmt.Println("Invalid option")
		}
	}
}
