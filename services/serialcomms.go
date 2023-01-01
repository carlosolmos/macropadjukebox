package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

const CMD_RESET = "reset"

type SerialComms struct {
	config *Configuration
	port   serial.Port
}

func NewSerialComms(config *Configuration) *SerialComms {
	sc := SerialComms{config: config}
	return &sc
}

func (sc *SerialComms) StartSerialComms() {
	if sc.config == nil {
		log.Fatal("OpenSerialComms: invalid configuration")
	}
	log.Println("start comms with serial port " + sc.config.SerialPort)

	found := false
	for {
		log.Println("list serial ports")
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			log.Fatal("No serial ports found!")
		}
		for _, portName := range ports {
			log.Printf("Port: %v\n", portName)
			if portName == sc.config.SerialPort {
				found = true
				break
			}
		}
		if found {
			break
		}
		time.Sleep(3 * time.Second)
	}

	log.Println("opening serial port")
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		DataBits: 8,
	}
	var err error
	sc.port, err = serial.Open(sc.config.SerialPort, mode)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.port.Close()
	err = sc.port.SetReadTimeout(serial.NoTimeout)
	if err != nil {
		log.Fatal(err)
	}

	// initial load of playlist
	sc.sendPlaylistCollection()

	log.Println("begin read cmd loop")
	buff := make([]byte, 256)
	for {
		// Reads bytes
		n, err := sc.port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			continue
		}
		input := string(buff[:n])
		if strings.LastIndex(strings.ToLower(input), "cmd|") == 0 {
			sc.processKeyCommand(input)
		} else if strings.LastIndex(strings.ToLower(input), "pl|") == 0 {
			sc.processPlaylistSelection(input)
		}
	}
}

func (sc *SerialComms) sendData(data []byte) error {
	_, err := sc.port.Write([]byte(fmt.Sprintf("%s\r", data)))
	if err != nil {
		return err
	}
	return nil
}

func (sc *SerialComms) sendPlaylistCollection() {
	log.Println("send collection to the macropad")
	collection := SimplerPlaylistMap(sc.config)
	jsonString, err := json.Marshal(collection)
	log.Println(string(jsonString))
	if err != nil {
		log.Println("error converting playlist map: " + err.Error())
	}
	err = sc.sendData(jsonString)
	if err != nil {
		log.Println("error sending playlist: " + err.Error())
	}
}

//Syntax: cmd|command_key
func (sc *SerialComms) processKeyCommand(cmd string) {
	cleanCommand := strings.ToLower(cmd)
	cleanCommand = strings.Trim(cleanCommand, " \r\n")
	cmdToks := strings.Split(cleanCommand, "|")
	if len(cmdToks) == 2 {
		log.Println("executing command " + cmdToks[1])

		if cmdToks[1] == CMD_RESET {
			sc.sendPlaylistCollection()
			err := VLC_Clear()
			if err != nil {
				log.Println(err.Error())
			}
			return
		}

		// key command
		cmdFun, ok := KeyFunctionsMap[cmdToks[1]]
		if ok {
			cmdFun()
		} else {
			log.Println(" command not found")
		}
	} else {
		log.Println(" invalid command string " + cleanCommand)
	}
}

func (sc *SerialComms) processPlaylistSelection(cmd string) {
	cleanCommand := strings.ToLower(cmd)
	cleanCommand = strings.Trim(cleanCommand, " \r\n")
	cmdToks := strings.Split(cleanCommand, "|")
	if len(cmdToks) == 2 {
		log.Println("play playlist " + cmdToks[1])
		if cmdToks[1] != "00" {
			pl, ok := sc.config.PlaylistCollection[cmdToks[1]]
			if ok {
				err := VLC_StartPlaylist(pl.File)
				if err != nil {
					log.Println("error playing playlist: " + err.Error())
				}
			} else {
				log.Println(" playlist not found")
			}
		}
	} else {
		log.Println(" invalid playlist string " + cleanCommand)
	}
}
