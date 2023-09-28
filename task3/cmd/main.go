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

	conn, err := ari.New("blondie", "192.168.0.16:8088", "asterisk", "test123")
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
			conn.Dial(args[1:]...)
		}
	}
}
