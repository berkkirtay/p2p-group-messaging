// Copyright (c) 2024 Berk Kirtay

package message

import (
	"context"
	"main/infra/store"
	"main/services/audit"
	"main/services/room"
	"slices"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
 * Basic messaging service on HTTP.
 *
 * This implementation uses permanent persistence of room instances.
 * Whereas in P2P implementation, a room can only be run by a peer.
 * The peers must be part of the room members to be able to host the
 * room and send messages.
 */

type MessageService interface {
	ReceiveMessages(id string, userId string, size string) []Message
	SendAMessage(id string, userId string, message string) Message
	DeleteAMessage(id string, userId string, message string) Message
}

var messageRepository = store.NewRepo("messaging")
var roomRepository = store.NewRepo("rooms")

func ReceiveMessages(id string, size string, sort string, userId string) []Message {
	var messages []Message = []Message{}
	var room room.Room = fetchTargetRoom(id)
	// Check if the user is in the room:
	if !validateUserRoomAuth(room, userId) {
		return messages
	}

	// Retrieve messages in the room:
	options := options.Find()
	var limit int64
	if size == "" {
		limit = 5
	} else {
		limit, _ = strconv.ParseInt(size, 10, 64)
	}
	options.SetLimit(limit)
	var sortValue int64 = -1
	if sort == "true" {
		sortValue = 1
	}
	options.SetSort(bson.M{"$natural": sortValue})
	filter := bson.D{{Key: "roomId", Value: id}}
	list, err := messageRepository.Find(filter, options)
	if err != nil && err != mongo.ErrNoDocuments {
		panic(err)
	} else {
		for list.Next(context.TODO()) {
			var currentMessage Message
			err := list.Decode(&currentMessage)
			if err != nil {
				panic(err)
			}
			updateOrDeleteMessageForRecipients(currentMessage, userId)
			messages = append(messages, currentMessage)
		}
	}

	return messages
}

func SendAMessage(id string, userId string, message Message) Message {
	// Check if the user is in the room:
	var room room.Room = fetchTargetRoom(id)
	if !validateUserRoomAuth(room, userId) {
		return CreateDefaultMessage()
	}

	// Build and send the message:
	var builtMessage Message = buildAMessage(room, userId, message)
	messageRepository.InsertOne(builtMessage)
	return builtMessage
}

func fetchTargetRoom(id string) room.Room {
	var room room.Room
	filter := bson.D{{Key: "id", Value: id}}
	cur, err := roomRepository.FindOne(filter, nil)
	if cur != nil && err == nil {
		cur.Decode(&room)
	}
	return room
}

func validateUserRoomAuth(room room.Room, userId string) bool {
	return slices.Contains(room.Members, userId)
}

// TODO change
func buildAMessage(room room.Room, userId string, message Message) Message {
	// Generate an id for message:
	var lastRecord Message = Message{}
	var newMessageId int
	options := options.FindOne().SetSort(bson.M{"$natural": -1})
	res, err := messageRepository.FindOne(bson.M{}, options)
	if res == nil && err == nil {
		// No message is found in the DB,
		// Generate a default id:
		newMessageId = 1000000000
	} else {
		res.Decode(&lastRecord)
		newMessageId, _ = strconv.Atoi(lastRecord.Id)
	}
	return CreateMessage(
		WithMessageId(strconv.Itoa(newMessageId+1)),
		WithUserId(userId),
		WithRoomId(room.Id),
		WithRecipients(room.Members),
		WithText(message.Text),
		WithIsEncrypted(message.IsEncrypted),
		WithMessageSignature(nil),
		WithMessageAudit(audit.CreateAuditForMessage()))
}

func updateOrDeleteMessageForRecipients(message Message, recipientToRemove string) {
	if recipientToRemove != "" {
		var index int = -1
		for i, recipient := range message.Recipients {
			if recipient == recipientToRemove {
				index = i
				break
			}
		}
		if index != -1 {
			message.Recipients = slices.Delete(message.Recipients, index, index+1)
		}
	}
	filter := bson.D{{Key: "id", Value: message.Id}}
	if len(message.Recipients) == 0 {
		//	messageRepository.DeleteOne(filter, nil)
	} else {
		messageRepository.ReplaceOne(filter, nil, message)
	}
}
