// Copyright (c) 2024 Berk Kirtay

package room

import (
	"main/infra/cryptography"
	"main/services/audit"
)

// Reminder: Owner of the room can be verified easily by checking the signature.

type Room struct {
	Id               string                     `json:"id,omitempty" bson:"id,omitempty"`
	Name             string                     `json:"name,omitempty" bson:"name,omitempty"`
	Info             string                     `json:"info,omitempty" bson:"info,omitempty"`
	Password         string                     `json:"password,omitempty" bson:"password,omitempty"`
	Capacity         int64                      `json:"capacity,omitempty" bson:"capacity,omitempty"`
	Members          []string                   `json:"members,omitempty" bson:"members,omitempty"`
	Signature        *cryptography.Cryptography `json:"signature,omitempty" bson:"signature,omitempty"`
	Audit            *audit.Audit               `json:"audit,omitempty" bson:"audit,omitempty"`
	DiffieHelmanKeys map[string]string          `json:"-"`
	RoomMasterKey    string                     `json:"roomMasterKey,omitempty" bson:"roomMasterKey,omitempty`
}

type RoomOption func(Room) Room

func WithId(id string) RoomOption {
	return func(room Room) Room {
		room.Id = id
		return room
	}
}

func WithName(name string) RoomOption {
	return func(room Room) Room {
		room.Name = name
		return room
	}
}

func WithInfo(info string) RoomOption {
	return func(room Room) Room {
		room.Info = info
		return room
	}
}

func WithPassword(password string) RoomOption {
	return func(room Room) Room {
		room.Password = password
		return room
	}
}

func WithCapacity(capacity int64) RoomOption {
	return func(room Room) Room {
		room.Capacity = capacity
		return room
	}
}

func WithMembers(members []string) RoomOption {
	return func(room Room) Room {
		room.Members = members
		return room
	}
}

func WithSignature(signature *cryptography.Cryptography) RoomOption {
	return func(room Room) Room {
		room.Signature = signature
		return room
	}
}

func WithAudit(audit *audit.Audit) RoomOption {
	return func(room Room) Room {
		room.Audit = audit
		return room
	}
}

func WithDiffieHelmanKeys(masterKeys map[string]string) RoomOption {
	return func(room Room) Room {
		room.DiffieHelmanKeys = masterKeys
		return room
	}
}

func WithRoomMasterKey(roomMasterKey string) RoomOption {
	return func(room Room) Room {
		room.RoomMasterKey = roomMasterKey
		return room
	}
}

func CreateDefaultRoom() Room {
	return Room{}
}

func CreateRoom(options ...RoomOption) Room {
	room := CreateDefaultRoom()

	for _, option := range options {
		room = option(room)
	}

	return room
}
