// Copyright (c) 2024 Berk Kirtay

package main

import (
	"main/p2p"
)

const (
	CLIENT_MODE = "--client_mode"
	CONNECT     = "--connect"
)

func main() {
	p2p.StartNode()
	p2p.StartClient()
}

// func parseCommandArguments(args []string) map[string]string {
// 	argumentPairs := make(map[string]string)
// 	for i, arg := range args {
// 		if i == 0 {
// 			continue
// 		}
// 		values := strings.Split(arg, "=")
// 		argumentPairs[values[0]] = values[1]
// 	}
// 	return argumentPairs
// }
