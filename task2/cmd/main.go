package main

import (
	amigo "ami/pkg"
	"log"
)

func main() {

	ami := amigo.NewManager()
	if err := ami.SetConf("10.1.0.228", "8088", "adel", "123"); err != nil {
		log.Fatalf(err.Error())
	}
	ami.Login()
}
