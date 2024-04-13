// Copyright (c) 2024 Berk Kirtay

package main

import (
	"main/p2p"
)

func main() {
	go p2p.StartNode()
	p2p.StartClient()
}
