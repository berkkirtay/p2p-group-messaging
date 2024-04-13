// Copyright (c) 2024 Berk Kirtay

package commands

import "main/services/peer"

var currentPeer peer.Peer
var api string

func InitializePeer(peer peer.Peer) {
	currentPeer = peer
	api = peer.Address
}

func GetCurrentPeer() peer.Peer {
	return currentPeer
}
