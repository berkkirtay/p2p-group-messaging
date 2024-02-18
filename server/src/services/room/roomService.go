package room

import (
	"context"
	"main/infrastructure"
	"main/services/audit"
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
	LeaveRoom(id string) Room                          // POST
	KickUser(roomId string, userId string, userToBeKicked string) bool
	BanUser(roomId string, userId string, userToBeBanned string) bool
}

var roomRepository = infrastructure.NewRepo("rooms")

func GetRooms(id string, size string) []Room {
	var rooms []Room = []Room{}

	if id != "" {
		var room Room
		filter := bson.D{{Key: "id", Value: id}}
		cur, err := roomRepository.FindOne(context.TODO(), filter, nil)
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
		list, err := roomRepository.Find(context.TODO(), bson.D{{}}, options)
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
	createdRoom := CreateRoom(
		WithId(room.Id),
		WithName(room.Name),
		WithInfo(room.Info),
		WithPassword(room.Password),
		WithCapacity(room.Capacity),
		WithMembers([]string{}),
		WithSignature(nil),
		WithAudit(audit.CreateAuditForRoom()))

	roomRepository.InsertOne(context.TODO(), createdRoom)
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
	cur, err := roomRepository.FindOne(context.TODO(), filter, nil)
	if cur != nil && err == nil {
		cur.Decode(&actualRoom)
	}

	// Check if the found room and the given room are equal for room auth and validation:
	if actualRoom.Password != room.Password || len(actualRoom.Members) >= int(actualRoom.Capacity) {
		return CreateDefaultRoom()
	}

	actualRoom.Members = append(actualRoom.Members, userId)
	actualRoom.Audit.LastOnlineDate = time.Now().Format(time.RFC1123)
	actualRoom.Audit.NumberOfActions += 1

	res, _ := roomRepository.ReplaceOne(filter, nil, actualRoom)
	if res.ModifiedCount == 0 {
		return CreateDefaultRoom()
	}
	return actualRoom
}

func LeaveRoom(id string) Room {
	return CreateDefaultRoom()
}
