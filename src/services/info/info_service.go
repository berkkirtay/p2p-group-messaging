// Copyright (c) 2024 Berk Kirtay

package info

import (
	"main/commands"
)

const (
	NEW_MESSAGE_INFO = "NEW_MESSAGE_INFO"
	ONLINE_INFO      = "ONLINE_INFO"
)

func HandlePeerInfo(info string) bool {
	if info == NEW_MESSAGE_INFO {
		//	fmt.Println("info received from " + commands.GetMainPeer().Hostname + ", own node: " + commands.GetNode().Hostname)
		if commands.GetNode().Address == commands.GetMainPeer().Address {
			commands.SendNotificationToPeers()
		}
		commands.NewMessageAlert()
	}
	return true
}
