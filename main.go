package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/carlosolmos/macropadjukebox/services"
)

var configuration services.Configuration

/*
Variables defined in configuration.go
KeyFunctionsMap
PlaylistSlots
*/

func main() {
	configContent, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("failed to read config file: ", err)
	}
	err = json.Unmarshal(configContent, &configuration)
	if err != nil {
		log.Fatal("failed to parse config json: ", err)
	}

	// set vlc host
	services.VlcHost = configuration.VlcHost

	// start serial comms
	serialComms := services.NewSerialComms(&configuration)
	serialComms.StartSerialComms()
}
