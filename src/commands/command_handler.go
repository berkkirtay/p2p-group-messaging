// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	AVAILABLE_COMMANDS string = "register\nlogin\nget-rooms\ncreate-room\njoin-room\nexit-room\nwhoami\nkick {user-name}"
	ROOM_MODE          string = "ROOM"
	COMMAND_MODE       string = "COMMAND"
)

var currentMode string

func InitializeCommandLine() {
	resp, err := http.Get(assignedPeer.Address + "/peer")
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic("Connection to a peer could not be initialized.")
	}
	fmt.Printf("A master peer is initialized: %s\n", assignedPeer.Hostname)
	fmt.Println("To see available commands, please type 'help'.")
	currentMode = COMMAND_MODE
}

func HandleInput() {
	reader := bufio.NewReader(os.Stdin)
	for {
		var sender string = ""
		if CurrentUser.Id != "" {
			sender = CurrentUser.Name
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
		HandleRegister(inputs)
	case "login":
		HandleLogin(inputs)
	case "get-rooms":
		HandleGetRooms()
	case "create-room":
		HandleCreateRoom(inputs)
	case "join-room":
		currentMode = ROOM_MODE
		HandleJoinRoom(inputs, CurrentUser)
	case "whoami":
		HandleWhoAmI()
	case "help":
		fmt.Println("Available commands are:")
		fmt.Println(AVAILABLE_COMMANDS)
	case "exit":
		DeletePeer(assignedPeer)
		os.Exit(0)
	default:
		fmt.Printf("Unddefined command '%s'\n", command)
	}
}

func handleRoom(command string) {
	inputs := strings.Split(command, " ")
	switch inputs[0] {
	case "kick":
		HandleKick(inputs)
	case "exit-room":
		HandleExitRoom()
		currentMode = COMMAND_MODE
	default:
		HandleText(command)
	}
}
