// Copyright (c) 2024 Berk Kirtay

package info

import "main/commands"

const (
	NEW_MESSAGE_INFO = "NEW_MESSAGE_INFO"
	ONLINE_INFO      = "ONLINE_INFO"
)

func HandlePeerInfo(info string) bool {
	if info == NEW_MESSAGE_INFO {
		commands.NewMessageAlert()
	}
	return true
}
