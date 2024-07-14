// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/http"
	"main/services/peer"
)

var assignedPeer peer.Peer

func InitializeAMasterPeer(hostname string, address string) {
	assignedPeer = peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.INBOUND))
	peer.PostPeer(assignedPeer)
}

func RegisterPeer(targetAddress string, hostname string, address string) {
	var newPeer peer.Peer = peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.OUTBOUND))
	body, err := json.Marshal(newPeer)
	if err != nil {
		panic("err")

	}
	res := http.POST(targetAddress+"/peer", string(body), &newPeer)
	if res.StatusCode != http.CREATED {
		fmt.Println(res)
		panic("err")
	}
	peer.PostPeer(peer.CreatePeer(
		peer.WithPeer(newPeer),
		peer.WithRole(peer.INBOUND)))
}

func DeletePeer(peer.Peer) {
	res := http.DELETE(assignedPeer.Address+"/peer", nil, "hostId", assignedPeer.Hostname)
	if res.StatusCode != http.OK {
		fmt.Printf("Error removing the peer.")
	}
}

func IsPeerInitialized() bool {
	var currentPeers []peer.Peer = peer.GetPeers()
	for _, currentPeer := range currentPeers {
		if currentPeer.Role == peer.OUTBOUND && isPeerOnline(currentPeer) {
			assignedPeer = currentPeer
		}
	}
	return assignedPeer.Address != ""
}

func isPeerOnline(peer peer.Peer) bool {
	res := http.GET(peer.Address+"/peer", nil)
	if res == nil || (res != nil && res.StatusCode != http.OK) {
		fmt.Printf("Peer %s is offline.\n", peer.Hostname)
		return false
	}
	return true
}
