package main

import (
	config "ari/configs"
	ari "ari/pkg"
	"log"
)

func main() {

	var err error
	ari := ari.New()

	ari.Cfg, err = config.LoadConfig("../configs/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	// ari.Cfg, err = config.MakeConfig("10.1.0.228", "8088", "asterisk", "test123")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	ari.Cfg.PrintCfg()
}
