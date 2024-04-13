// Copyright (c) 2024 Berk Kirtay

package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"main/commands"
	"main/services/peer"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const AVAILABLE_COMMANDS = "register\nlogin\nget-rooms\ncreate-room\njoin-room\nexit-room\nwhoami\nkick {user-name}"

// TODO: DB sync between peers is neecessary.

var currentMode string

func StartClient() {
	time.Sleep(100 * time.Millisecond)
	handlePeerConnection()
	initializeCommandLine()
	handleInput()
}

func handlePeerConnection() {
	var localAddresses []string = []string{}
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if addr.To4() != nil {
			fmt.Println(addr.To4().String())
			localAddresses = append(localAddresses, addr.To4().String())
		}
	}
	for _, host := range localAddresses {
		var peers []peer.Peer = []peer.Peer{}
		resp, err := http.Get("http://" + host + ADDRESS + API + "peer")
		if err != nil {
			fmt.Println(err)
		}
		err = json.NewDecoder(resp.Body).Decode(&peers)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, peer := range peers {
			if peer.Role == "ACTIVE" {
				commands.InitializePeer(peer)
				goto done
			}
		}
	}

	if commands.GetCurrentPeer().Address == "" {
		fmt.Println("No active peer found, making yourself an active peer.")
		newPeer := peer.CreatePeer(peer.WithAddress("http://127.0.0.1:8080"), peer.WithRole("ACTIVE"))
		peer.PostPeer(newPeer)
		commands.InitializePeer(newPeer)
	}
done:
	fmt.Printf("A master peer is initialized: %s\n", commands.GetCurrentPeer())
}

func initializeCommandLine() {
	resp, err := http.Get(commands.GetCurrentPeer().Address + "/api/peer")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != OK {
		panic("Connection to server could not initialized.")
	}
	fmt.Println("Connection to server is initialized.")
	fmt.Println("To see available commands, please type 'help'.")

	currentMode = COMMAND_MODE
}

func handleInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var sender string = ""
		if commands.CurrentUser.Id != "" {
			sender = commands.CurrentUser.Name
		}
		fmt.Printf("%s >> ", sender)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		if len(text) == 0 {
			continue
		}
		if currentMode == COMMAND_MODE {
			handleCommand(text)
		} else {
			handleRoom(text)
		}
	}

}

func handleCommand(command string) {
	inputs := strings.Split(command, " ")
	switch inputs[0] {
	case "register":
		commands.HandleRegister(inputs)
	case "login":
		commands.HandleLogin(inputs)
	case "get-rooms":
		commands.HandleGetRooms()
	case "create-room":
		commands.HandleCreateRoom(inputs)
	case "join-room":
		currentMode = ROOM_MODE
		commands.HandleJoinRoom(inputs, commands.CurrentUser)
	case "whoami":
		commands.HandleWhoAmI()
	case "help":
		fmt.Println("Available commands are:")
		fmt.Println(AVAILABLE_COMMANDS)
	default:
		fmt.Printf("Unddefined command '%s'\n", command)
	}
}

func handleRoom(command string) {
	inputs := strings.Split(command, " ")
	switch inputs[0] {
	case "kick":
		commands.HandleKick(inputs)
	case "exit-room":
		commands.HandleExitRoom()
		currentMode = COMMAND_MODE
	default:
		commands.HandleText(command)
	}
}
