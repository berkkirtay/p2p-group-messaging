// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/cryptography"
	"main/infra/http"
	"main/services/peer"
)

var assignedPeer peer.Peer

func InitializeAMasterPeer(hostname string, address string) {
	assignedPeer = peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.INBOUND),
		peer.WithCryptography(cryptography.CreateCommonCrypto()))
	peer.PostPeer(assignedPeer)
}

// TODO review here:
func RegisterPeer(targetAddress string, hostname string, address string) {
	//peer.UpdatePeerEllipticKeys(&assignedPeer)
	var newPeer peer.Peer = peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.OUTBOUND),
		peer.WithCryptography(assignedPeer.Cryptography))
	body, err := json.Marshal(newPeer)
	if err != nil {
		panic("err")
	}
	res := http.POST(assignedPeer, targetAddress+"/peer", string(body), &newPeer)
	if res.StatusCode != http.CREATED {
		panic("err")
	}
	peer.PostPeer(peer.CreatePeer(
		peer.WithPeer(newPeer),
		peer.WithRole(peer.INBOUND),
		peer.WithCryptography(
			cryptography.CreateCryptography(
				cryptography.WithElliptic(assignedPeer.Cryptography.Elliptic)))))
}

func DeletePeer(peer.Peer) {
	res := http.DELETE(assignedPeer, assignedPeer.Address+"/peer", nil, "hostId", assignedPeer.Hostname)
	if res.StatusCode != http.OK {
		fmt.Printf("Error removing the peer.\n")
	}
}

func IsPeerInitialized() bool {
	var currentPeers []peer.Peer = peer.GetPeers(string(""), string(""))
	for _, currentPeer := range currentPeers {
		if currentPeer.Role == peer.OUTBOUND {
			assignedPeer = synchronizeWithPeer(currentPeer)
		}
	}
	return assignedPeer.Address != ""
}

func synchronizeWithPeer(currentPeer peer.Peer) peer.Peer {
	synchronizedPeer := make([]peer.Peer, 5)
	res := http.GET(
		currentPeer,
		currentPeer.Address+"/peer",
		&synchronizedPeer,
		"hostname",
		currentPeer.Hostname)
	if res == nil || res.StatusCode != http.OK {
		fmt.Printf("Peer %s is offline.\n", currentPeer.Hostname)
		return peer.CreateDefaultPeer()
	}
	return synchronizedPeer[0]
}
