package services

import (
	"fmt"
	"log"
)

/*
from the json file
"playlists": {
    "01": {
      "name": "Death Cab For Cutie - Transatlanticism",
      "file": "/home/pi/Music/DeathCabForCutie-Transatlanticism.m3u"
    },
	...
*/

type Playlist struct {
	Name string `json:"name"`
	File string `json:"file,omitempty"`
}

type Configuration struct {
	VlcHost            string              `json:"vlcHost"`
	SerialPort         string              `json:"serialPort"`
	PlaylistCollection map[string]Playlist `json:"playlists"`
}

func SimplerPlaylistMap(config *Configuration) map[string]string {
	output := map[string]string{}
	for k, pl := range config.PlaylistCollection {
		output[k] = pl.Name
	}
	return output
}

var KeyFunctionsMap = map[string]func(){
	"play": func() {
		fmt.Println("play")
		err := VLC_Play()
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"pause": func() {
		fmt.Println("pause")
		err := VLC_Pause()
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"stop": func() {
		fmt.Println("stop")
		err := VLC_Stop()
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"prev": func() {
		fmt.Println("prev")
		err := VLC_Prev()
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"next": func() {
		fmt.Println("next")
		err := VLC_Next()
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"volup": func() {
		fmt.Println("volup")
		err := VLC_Volup(1)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"voldown": func() {
		fmt.Println("voldown")
		err := VLC_Voldown(1)
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
	"clear": func() {
		fmt.Println("clear")
		err := VLC_Clear()
		if err != nil {
			log.Printf("error: %s\n", err.Error())
		}
	},
}
