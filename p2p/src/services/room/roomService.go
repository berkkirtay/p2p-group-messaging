package room

import (
	"context"
	"main/infrastructure"
	"main/services/audit"
	"main/services/cryptography"
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

var repository = infrastructure.NewRepo("rooms")

func GetRooms(id string, size string) []Room {
	var rooms []Room = []Room{}

	if id != "" {
		var room Room
		filter := bson.D{{Key: "id", Value: id}}
		cur, err := repository.FindOne(context.TODO(), filter, nil)
		if cur != nil && err == nil {
			cur.Decode(&room)
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
		list, err := repository.Find(context.TODO(), bson.D{{}}, options)
		if err != nil && err != mongo.ErrNoDocuments {
			panic(err)
		} else {
			for list.Next(context.TODO()) {
				var currentRoom Room
				err := list.Decode(&currentRoom)
				if err != nil {
					panic(err)
				}
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
	res, err := repository.FindOne(context.TODO(), bson.M{}, options)
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
		WithSignature(nil),
		WithAudit(audit.CreateAuditForRoom()))

	repository.InsertOne(context.TODO(), createdRoom)
	return createdRoom
}

func DeleteRooms(ids []string) []Room {
	return nil
}

func UpdateRoom(id string, room Room) Room {
	return CreateDefaultRoom()
}

func JoinRoom(id string, room Room, userId string) Room {
	var actualRoom Room
	filter := bson.D{{Key: "id", Value: id}}
	cur, err := repository.FindOne(context.TODO(), filter, nil)
	if cur != nil && err == nil {
		cur.Decode(&actualRoom)
	}

	if slices.Contains(actualRoom.Members, userId) {
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

	// masterSecret := cryptography.ServerSideDiffieHelmanKeyExhange(userHandshakeKey)
	// actualRoom.DiffieHelmanKeys[userId] = masterSecret[0]
	actualRoom.HandshakeKey = actualRoom.RoomMasterKey

	actualRoom.Members = append(actualRoom.Members, userId)
	actualRoom.Audit.LastOnlineDate = time.Now().Format(time.RFC1123)
	actualRoom.Audit.NumberOfActions += 1

	res, _ := repository.ReplaceOne(filter, nil, actualRoom)
	if res.ModifiedCount == 0 {
		return CreateDefaultRoom()
	}
	//	SendAMessage(id, userId, buildAMessage(room, userId, CreateMessage(WithText("Greetings! I just joined."))))
	return actualRoom
}

func LeaveRoom(id string, userId string) bool {
	var actualRoom Room
	filter := bson.D{{Key: "id", Value: id}}
	cur, err := repository.FindOne(context.TODO(), filter, nil)
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
