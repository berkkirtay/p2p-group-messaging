// Copyright (c) 2024 Berk Kirtay

package room

import (
	"context"
	"main/infra/cryptography"
	"main/infra/store"
	"main/services/audit"
	"main/services/user"
	"slices"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RoomService interface {
	GetRooms(id string, size string) []Room            // GET
	CreateRooms(rooms []Room) []Room                   // POST
	DeleteRooms(ids []string) []Room                   // DELETE
	UpdateRoom(id string, room Room) Room              // PUT
	JoinRoom(id string, room Room, userId string) Room // POST
	LeaveRoom(id string) bool                          // POST
	KickUser(roomId string, userId string, userToBeKicked string) bool
	BanUser(roomId string, userId string, userToBeBanned string) bool
}

var repository = store.NewRepo("rooms")

func GetRooms(id string, size string) []Room {
	var rooms []Room = []Room{}

	if id != "" {
		var room Room
		filter := bson.D{{Key: "id", Value: id}}
		cur, err := repository.FindOne(filter, nil)
		if cur != nil && err == nil {
			cur.Decode(&room)
			room.Password = ""
			room.RoomMasterKey = ""
			rooms = append(rooms, room)
		}
	} else {
		options := options.Find()
		var limit int64
		if size == "" {
			limit = 5
		} else {
			limit, _ = strconv.ParseInt(size, 10, 64)
		}
		options.SetLimit(limit)
		options.SetSort(bson.M{"$natural": -1})
		list, err := repository.Find(bson.D{{}}, options)
		if err != nil && err != mongo.ErrNoDocuments {
			panic(err)
		} else {
			for list.Next(context.TODO()) {
				var currentRoom Room
				err := list.Decode(&currentRoom)
				if err != nil {
					panic(err)
				}
				currentRoom.Password = ""
				currentRoom.RoomMasterKey = ""
				rooms = append(rooms, currentRoom)
			}
		}
	}
	return rooms
}

func PostRoom(room Room) Room {
	var newRoomId int
	var lastRecord Room = Room{}
	options := options.FindOne().SetSort(bson.M{"$natural": -1})
	res, err := repository.FindOne(bson.M{}, options)
	if res == nil && err == nil {
		newRoomId = 12345
	} else {
		res.Decode(&lastRecord)
		newRoomId, _ = strconv.Atoi(lastRecord.Id)
	}
	createdRoom := CreateRoom(
		WithId(strconv.Itoa(newRoomId+1)),
		WithName(room.Name),
		WithInfo(room.Info),
		WithPassword(room.Password),
		WithCapacity(room.Capacity),
		WithMembers([]string{}),
		WithRoomMasterKey(cryptography.GenerateARandomMasterSecret()),
		WithSignature(cryptography.CreateCommonCrypto(
			room.Name,
			room.Info,
			room.Password,
			room.RoomMasterKey,
		)),
		WithAudit(audit.CreateAuditForRoom()))

	repository.InsertOne(createdRoom)
	return createdRoom
}

func DeleteRooms(ids []string) int64 {
	var deletedCount int64 = 0
	for _, id := range ids {
		filter := bson.D{{Key: "id", Value: id}}
		res, _ := repository.DeleteOne(filter, nil)
		deletedCount += res.DeletedCount
	}
	return deletedCount
}

func UpdateRoom(id string, room Room) Room {
	return CreateDefaultRoom()
}

func JoinRoom(id string, room Room, userId string) Room {
	var actualRoom Room
	filter := bson.D{{Key: "id", Value: id}}
	cur, err := repository.FindOne(filter, nil)
	if cur != nil && err == nil {
		cur.Decode(&actualRoom)
	}

	user := user.GetUser(userId, "")

	if slices.Contains(actualRoom.Members, userId) {
		actualRoom.RoomMasterKey = cryptography.EncryptRSA(
			actualRoom.RoomMasterKey,
			user.Cryptography.PublicKey)
		return actualRoom
	}

	/*
	 * Check if the found room and the given room
	 * are equal for room authentication and validation:
	 */
	if actualRoom.Password != room.Password ||
		len(actualRoom.Members) >= int(actualRoom.Capacity) {
		return CreateDefaultRoom()
	}

	//actualRoom.RoomMasterKey = cryptography.EnrichMasterSecret(actualRoom.RoomMasterKey, user.Signature.Hash)
	actualRoom.Members = append(actualRoom.Members, userId)
	actualRoom.Audit.LastOnlineDate = time.Now().Format(time.RFC1123)
	actualRoom.Audit.NumberOfActions += 1

	res, _ := repository.ReplaceOne(filter, nil, actualRoom)
	if res.ModifiedCount == 0 {
		return CreateDefaultRoom()
	}
	//	SendAMessage(id, userId, buildAMessage(room, userId, CreateMessage(WithText("Greetings! I just joined."))))
	actualRoom.RoomMasterKey = cryptography.EncryptRSA(actualRoom.RoomMasterKey, user.Cryptography.PublicKey)
	return actualRoom
}

func LeaveRoom(id string, userId string) bool {
	var actualRoom Room
	filter := bson.D{{Key: "id", Value: id}}
	cur, err := repository.FindOne(filter, nil)
	if cur != nil && err == nil {
		cur.Decode(&actualRoom)
	}

	var index int = -1
	for i, roomUserId := range actualRoom.Members {
		if roomUserId == userId {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}
	actualRoom.Members = slices.Delete(actualRoom.Members, index, index+1)
	res, _ := repository.ReplaceOne(filter, nil, actualRoom)
	if res.ModifiedCount == 0 {
		return false
	}

	return true
}
