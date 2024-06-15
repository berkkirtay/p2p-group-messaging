// Copyright (c) 2024 Berk Kirtay

package p2p

import (
	"main/commands"
)

// TODO: DB sync between peers is neecessary.

func StartClient() {
	commands.InitializeCommandLine()
	commands.HandleInput()
}
