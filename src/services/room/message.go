// Copyright (c) 2024 Berk Kirtay

package room

import (
	"main/services/audit"
	"main/services/cryptography"
)

type Message struct {
	Id          string                  `json:"id,omitempty" bson:"id,omitempty"`
	UserId      string                  `json:"userId,omitempty" bson:"userId,omitempty"`
	RoomId      string                  `json:"roomId,omitempty" bson:"roomId,omitempty"`
	Text        string                  `json:"text,omitempty" bson:"text,omitempty"`
	Signature   *cryptography.Signature `json:"signature,omitempty" bson:"signature,omitempty"`
	IsEncrypted bool                    `json:"isEncrypted,omitempty" bson:"isEncrypted,omitempty"`
	Audit       *audit.Audit            `json:"audit,omitempty" bson:"audit,omitempty"`
}

type MessageOption func(Message) Message

func WithMessageId(id string) MessageOption {
	return func(message Message) Message {
		message.Id = id
		return message
	}
}

func WithUserId(userId string) MessageOption {
	return func(message Message) Message {
		message.UserId = userId
		return message
	}
}

func WithRoomId(roomId string) MessageOption {
	return func(message Message) Message {
		message.RoomId = roomId
		return message
	}
}

func WithText(text string) MessageOption {
	return func(message Message) Message {
		message.Text = text
		return message
	}
}

func WithMessageSignature(signature *cryptography.Signature) MessageOption {
	return func(message Message) Message {
		message.Signature = signature
		return message
	}
}

func WithIsEncrypted(isEncrypted bool) MessageOption {
	return func(message Message) Message {
		message.IsEncrypted = isEncrypted
		return message
	}
}

func WithMessageAudit(audit *audit.Audit) MessageOption {
	return func(message Message) Message {
		message.Audit = audit
		return message
	}
}

func CreateDefaultMessage() Message {
	return Message{}
}

func CreateMessage(options ...MessageOption) Message {
	message := CreateDefaultMessage()

	for _, option := range options {
		message = option(message)
	}

	return message
}
