// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/cryptography"
	"main/infra/http"
	"main/services/room"
	"main/services/user"
	"strconv"
	"time"
)

var currentUser user.User
var currentRoom room.Room
var roomUsers map[string]user.User
var retrieveTextsFlag bool = true
var lastId string

func HandleGetRooms() {
	url := assignedPeer.Address + "/room"
	var rooms = make([]room.Room, 5)
	var res = http.GET(url, &rooms, "size", "5")
	if res.StatusCode != 200 {
		fmt.Printf("Error")
		return
	}
	fmt.Printf("Available rooms in the server:\n")
	fmt.Printf("------------\n")
	for _, room := range rooms {
		fmt.Printf("Id and Room Name: %s - %s\nInfo: %s\n", room.Id, room.Name, room.Info)
		fmt.Printf("Capacity: %v\nOther details: %s\n", room.Capacity, room.Audit.CreateDate)
		fmt.Printf("------------\n")
	}
}

// TODO: changing the input handling way to 'create-room --room-name={room-name}...'
func HandleCreateRoom(command []string) {
	if len(command) < 4 {
		fmt.Printf("Wrong usage.\n(create-room {room-name} {info} {capacity} {password})\n")
		return
	}
	capacity, err := strconv.ParseInt(command[3], 10, 64)
	if err != nil {
		return
	}
	var room = room.CreateRoom(room.WithName(command[1]),
		room.WithInfo(command[2]),
		room.WithCapacity(capacity),
		room.WithPassword(command[4]))
	body, err := json.Marshal(room)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	res := http.POST(assignedPeer.Address+"/room", string(body), &room)
	if res.StatusCode != 201 {
		fmt.Printf("Error")
		return
	}
	fmt.Printf("Room is created successfully with the id: %s\n", room.Id)
}

func HandleJoinRoom(command []string, user user.User) {
	currentUser = user
	joinRoom(command[1], command[2])
	go retrieveTexts()
	retrieveTextsFlag = true
}

func joinRoom(roomId string, roomPassword string) {
	url := assignedPeer.Address + "/room/join"
	var room = room.CreateRoom(
		room.WithId(roomId),
		room.WithPassword(roomPassword))
	body, err := json.Marshal(room)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	res := http.POST(url, string(body), &room, "id", roomId)
	if res.StatusCode != 200 {
		fmt.Printf("Error")
		return
	}
	currentRoom = room
	fmt.Printf("Joined the room. You will talk with:\n")
	roomUsers = make(map[string]user.User)
	for _, userId := range room.Members {
		var userBody = []user.User{}
		var res = http.GET(assignedPeer.Address+"/users", &userBody, "id", userId)
		if res.StatusCode != 200 {
			fmt.Printf("Error")
			return
		}
		roomUsers[userId] = userBody[0]
		fmt.Printf("%s\n", userBody[0].Name)
	}
}

func HandleText(command string) {
	url := assignedPeer.Address + "/room/messages"
	var message room.Message = room.CreateMessage(
		room.WithText(cryptography.Encrypt(command, currentRoom.HandshakeKey)),
		room.WithIsEncrypted(true))
	body, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	res := http.POST(url, string(body), message, "id", currentRoom.Id)
	if res.StatusCode != 201 {
		fmt.Printf("Error")
		return
	}
}

func retrieveTexts() {
	for {
		if retrieveTextsFlag {
			time.Sleep(90 * time.Millisecond)
			getTexts()
		} else {
			break
		}
	}
}

func getTexts() {
	url := assignedPeer.Address + "/room/messages"
	var messages = []room.Message{}
	res := http.GET(url, &messages, "id", currentRoom.Id, "size", "1")
	if res.StatusCode == 200 && len(messages) > 0 &&
		lastId != messages[0].Id && messages[0].UserId != currentUser.Id {
		fmt.Printf("\r%s >> %s\n", roomUsers[messages[0].UserId].Name,
			buildAReadableText(messages[0]))
		fmt.Printf("%s >> ", currentUser.Name)
		lastId = messages[0].Id
	}
}

func buildAReadableText(message room.Message) string {
	if message.IsEncrypted {
		return cryptography.Decrypt(message.Text, currentRoom.HandshakeKey)
	} else {
		return message.Text
	}
}

func HandleExitRoom() {
	fmt.Println("You left the room.")
}

func HandleKick(command []string) {}
