// Copyright (c) 2024 Berk Kirtay

package message

import (
	"main/infra/cryptography"
	"main/services/audit"
)

type Message struct {
	Id          string                     `json:"id,omitempty" bson:"id,omitempty"`
	UserId      string                     `json:"userId,omitempty" bson:"userId,omitempty"`
	RoomId      string                     `json:"roomId,omitempty" bson:"roomId,omitempty"`
	Recipients  []string                   `json:"recipients,omitempty" bson:"recipients,omitempty"`
	Text        string                     `json:"text,omitempty" bson:"text,omitempty"`
	Signature   *cryptography.Cryptography `json:"signature,omitempty" bson:"signature,omitempty"`
	IsEncrypted bool                       `json:"isEncrypted,omitempty" bson:"isEncrypted,omitempty"`
	Audit       *audit.Audit               `json:"audit,omitempty" bson:"audit,omitempty"`
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

func WithRecipients(recipients []string) MessageOption {
	return func(message Message) Message {
		message.Recipients = recipients
		return message
	}
}

func WithText(text string) MessageOption {
	return func(message Message) Message {
		message.Text = text
		return message
	}
}

func WithMessageSignature(signature *cryptography.Cryptography) MessageOption {
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

func WithMessage(newMessage Message) MessageOption {
	return func(message Message) Message {
		message.Id = newMessage.Id
		message.UserId = newMessage.UserId
		message.RoomId = newMessage.RoomId
		message.Recipients = newMessage.Recipients
		message.Text = newMessage.Text
		message.Signature = newMessage.Signature
		message.IsEncrypted = newMessage.IsEncrypted
		message.Audit = newMessage.Audit
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
