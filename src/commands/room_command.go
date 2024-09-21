// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/cryptography"
	"main/infra/http"
	"main/services/room"
	"main/services/user"
	"sort"
	"strconv"
	"time"
)

// TODO: To improve performance, use message numbers for rooms.

var currentUser user.User
var currentRoom room.Room
var roomUsers map[string]user.User
var retrieveMessagesFlag bool
var lastMessageId int64

func HandleGetRooms() {
	var rooms = make([]room.Room, 5)
	var res = http.GET(assignedPeer, assignedPeer.Address+"/room", &rooms, "size", "5")
	if res.StatusCode == http.OK {
		fmt.Printf("Available rooms in the server:\n")
		fmt.Printf("------------\n")
		for _, room := range rooms {
			fmt.Printf("Id and Room Name: %s - %s\nInfo: %s\n", room.Id, room.Name, room.Info)
			fmt.Printf("Capacity: %v\nOther details: %s\n", room.Capacity, room.Audit.CreateDate)
			fmt.Printf("------------\n")
		}
	} else {
		fmt.Printf("No rooms found.")
	}
}

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
	res := http.POST(assignedPeer, assignedPeer.Address+"/room", string(body), &room)
	if res.StatusCode != http.CREATED {
		fmt.Printf("Error")
		return
	}
	fmt.Printf("Room is created successfully with the id: %s\n", room.Id)
}

func HandleText(command string) {
	var message room.Message = room.CreateMessage(
		room.WithText(cryptography.EncryptAES(command, currentRoom.RoomMasterKey)),
		room.WithIsEncrypted(true))
	body, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	res := http.POST(assignedPeer, assignedPeer.Address+"/room/messages", string(body), message, "id", currentRoom.Id)
	if res.StatusCode != http.CREATED {
		fmt.Printf("Message could not be sent.")
		return
	}
	fmt.Printf("\033[1A\033[K")
}

func HandleJoinRoom(command []string, user user.User) {
	currentUser = user
	joinRoom(command[1], command[2])
	retrieveMessagesFlag = true
	go messageLoop()
}

func joinRoom(roomId string, roomPassword string) {
	var room = room.CreateRoom(
		room.WithId(roomId),
		room.WithPassword(roomPassword))
	body, err := json.Marshal(room)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	res := http.POST(assignedPeer, assignedPeer.Address+"/room/join", string(body), &room, "id", roomId)
	if res.StatusCode != http.OK {
		fmt.Printf("Error")
		return
	}
	currentRoom = room
	currentRoom.RoomMasterKey = cryptography.DecryptRSA(
		currentRoom.RoomMasterKey, sessionAuth.Cryptography.PrivateKey)

	fmt.Printf("Joined the room. You will talk with:\n")
	roomUsers = make(map[string]user.User)
	for _, userId := range room.Members {
		roomUsers[userId] = getUser(userId)
		fmt.Printf("%s\n", roomUsers[userId].Name)
	}
}

func messageLoop() {
	for {
		if !retrieveMessagesFlag {
			break
		}
		var messageSize int64 = fetchLastMessageId() - lastMessageId
		if messageSize > 0 {
			getMessages(messageSize)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func getMessages(size int64) {
	var messages = []room.Message{}
	res := http.GET(
		assignedPeer,
		assignedPeer.Address+"/room/messages",
		&messages,
		"id",
		currentRoom.Id,
		"size",
		strconv.FormatInt(size, 10))
	if res.StatusCode != http.INTERNAL_SERVER_ERROR {
		printMessages(messages)
		lastMessageId, _ = strconv.ParseInt(messages[len(messages)-1].Id, 10, 64)
	}
}

func printMessages(messages []room.Message) {
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Id < messages[j].Id
	})

	for _, message := range messages {
		currentMessageId, _ := strconv.ParseInt(message.Id, 10, 64)
		if currentMessageId > lastMessageId {
			if roomUsers[message.UserId].Id == "" {
				roomUsers[message.UserId] = getUser(message.UserId)
			}
			fmt.Printf(
				"\r%s >> %s\n", roomUsers[message.UserId].Name,
				buildAReadableText(message))
			fmt.Printf("%s >> ", currentUser.Name)
		}
	}
}

func fetchLastMessageId() int64 {
	var messages = []room.Message{}
	res := http.GET(
		assignedPeer,
		assignedPeer.Address+"/room/messages",
		&messages,
		"id",
		currentRoom.Id,
		"size",
		"1")
	if res != nil && res.StatusCode != http.NOT_FOUND {
		lastStoredMessageId, _ := strconv.ParseInt(messages[0].Id, 10, 64)
		return lastStoredMessageId
	}
	return 0
}

func buildAReadableText(message room.Message) string {
	if message.IsEncrypted {
		return cryptography.DecryptAES(message.Text, currentRoom.RoomMasterKey)
	} else {
		return message.Text
	}
}

func getUser(userId string) user.User {
	var userBody = []user.User{}
	var res = http.GET(
		assignedPeer,
		assignedPeer.Address+"/users",
		&userBody,
		"id",
		userId)
	if res.StatusCode == http.OK {
		return userBody[0]
	}
	return user.CreateDefaultUser()
}

func HandleExitRoom() {
	retrieveMessagesFlag = false
	fmt.Println("You left the room.")
}

func HandleKick(command []string) {}
