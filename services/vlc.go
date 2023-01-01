package services

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

/**
VLC remote control via the socket interface.

VLC has to be running with the RC interface enabled
https://wiki.videolan.org/Documentation:Advanced_Use_of_VLC/#RC_and_RTCI

headless:
cvlc --rc-host="localhost:8888" -I rc
with GUI:
vlc --rc-host="localhost:8888" --extraintf rc

The interface is interactive and responds to string commands.

List of commands: see vlc.md

*/

var VlcHost = "192.168.86.176:8888"

func VLCsendCommand(cmdList []string, waitForResponse bool) (string, error) {
	log.Printf("sending commands %s\n", cmdList)
	//return "", nil
	vlcsocket, err := net.Dial("tcp", VlcHost)
	if err != nil {
		log.Println("Connection to socket failed:", err)
		return "", err
	}
	// wait for the prompt to respond.
	time.Sleep(1 * time.Second)
	// read prompt out of the way
	recvBuf := make([]byte, 1024)
	n, err := vlcsocket.Read(recvBuf[:]) // recv data

	// send commands
	res := ""
	for _, cmd := range cmdList {
		cmdString := fmt.Sprintf("%s\n", cmd)
		_, err = vlcsocket.Write([]byte(cmdString))
		if err != nil {
			log.Println("Write to socket failed:", err)
			return "", err
		}
		// read command response
		if waitForResponse {
			err = vlcsocket.SetReadDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				log.Println("SetReadDeadline failed:", err)
				// do something else, for example create new conn
				return "", err
			}
			// try 3 times
			for i := 0; i < 3; i++ {
				recvBuf := make([]byte, 1024)
				n, err = vlcsocket.Read(recvBuf[:]) // recv data
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						//timeout
						break
					} else {
						// some error else, do something else, for example create new conn
						log.Println("Read from socket error:", err)
						return res, nil
					}
				}
				if n == 0 {
					break
				}
				res = res + string(recvBuf)
			}
		}
	}
	// close the connection gracefully
	_, err = vlcsocket.Write([]byte("quit\n"))
	if err != nil {
		return res, err
	}
	_ = vlcsocket.Close()
	return res, nil
}

func VLC_Status() (string, error) {
	resp, err := VLCsendCommand([]string{"status"}, true)
	if err != nil {
		return "", err
	}
	resp = strings.Replace(resp, ">", "", -1)
	resp = strings.TrimSpace(resp)
	return resp, nil
}

func VLC_IsPlaying() (bool, error) {
	resp, err := VLCsendCommand([]string{"is_playing"}, true)
	if err != nil {
		return false, err
	}
	resp = strings.Replace(resp, ">", "", -1)
	resp = strings.TrimSpace(resp)
	return resp == "1", err
}

func VLC_StartPlaylist(playlist string) error {
	cmdList := []string{
		"clear",
		fmt.Sprintf("add %s", playlist),
	}
	_, err := VLCsendCommand(cmdList, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Add(xyz string) error {
	_, err := VLCsendCommand([]string{fmt.Sprintf("add %s", xyz)}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Enqueue(xyz string) error {
	_, err := VLCsendCommand([]string{fmt.Sprintf("queue %s", xyz)}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Clear() error {
	_, err := VLCsendCommand([]string{"clear"}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Play() error {
	_, err := VLCsendCommand([]string{"play"}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Pause() error {
	_, err := VLCsendCommand([]string{"pause"}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Stop() error {
	_, err := VLCsendCommand([]string{"stop"}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Next() error {
	_, err := VLCsendCommand([]string{"next"}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Prev() error {
	_, err := VLCsendCommand([]string{"prev"}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Volup(step int) error {
	_, err := VLCsendCommand([]string{fmt.Sprintf("volup %d", step)}, false)
	if err != nil {
		return err
	}
	return nil
}

func VLC_Voldown(step int) error {
	_, err := VLCsendCommand([]string{fmt.Sprintf("voldown %d", step)}, false)
	if err != nil {
		return err
	}
	return nil
}
