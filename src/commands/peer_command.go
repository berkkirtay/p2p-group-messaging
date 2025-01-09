// Copyright (c) 2024 Berk Kirtay

package commands

import (
	"encoding/json"
	"fmt"
	"main/infra/cryptography"
	"main/infra/http"
	"main/infra/network"
	"main/services/peer"
)

var assignedPeer peer.Peer

func InitializePeer() {
	network.InitializeBroadcast()
	network.StartPeerBroadcast()
	go peerListenerLoop()

	var remotePeers []peer.Peer = fetchAvailablePeers()
	assignedPeer = handleUserPeerSelection(remotePeers)
	peer.PostPeer(assignedPeer)
	fmt.Printf("A master peer is initialized: %s\n", assignedPeer.Hostname)
}

func DeletePeer(peer.Peer) {
	res := http.DELETE(
		assignedPeer,
		assignedPeer.Address+"/peer",
		nil,
		"hostId",
		assignedPeer.Hostname)
	if res.StatusCode != http.OK {
		fmt.Printf("Error removing the peer.\n")
	}
}

func peerListenerLoop() {
	for {
		var targetAddress string = network.ListenForPeerBroadcast()
		hostname, address := network.GetHostAddress()
		registerPeer(targetAddress, hostname, address)
	}
}

func registerPeer(targetAddress string, hostname string, address string) {
	var newPeer peer.Peer = peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.OUTBOUND),
		peer.WithCryptography(cryptography.CreateCommonCrypto()))

	body, err := json.Marshal(newPeer)
	if err != nil {
		panic("err")
	}

	// Sending own node information to the incoming remote peer
	res := http.POST(assignedPeer, targetAddress+"/peer", string(body), &newPeer)
	if res.StatusCode != http.CREATED {
		panic("err")
	}
	// Saving the remote peer information to your own database
	peer.PostPeer(peer.CreatePeer(
		peer.WithPeer(newPeer),
		peer.WithRole(peer.INBOUND)))
}

func handleUserPeerSelection(peers []peer.Peer) peer.Peer {
	if len(peers) > 0 {
		fmt.Println("--------------")
		fmt.Println("i |	Hostname	|	Address		")
		for i, peer := range peers {
			fmt.Printf("%d |	%s	|	%s		\n", i, peer.Hostname, peer.Address)
		}
		fmt.Println("--------------")

		fmt.Println("Please select a peer to make a connection (-1 for default): ")
		var number int
		_, err := fmt.Scanf("%d", &number)
		if err == nil && number <= len(peers)-1 && number != -1 {
			return peers[number]
		}
	}
	fmt.Println("No connection has been done to any peer. Initializing this node as an active peer.")
	return initializeAMasterPeer()
}

func initializeAMasterPeer() peer.Peer {
	hostname, address := network.GetHostAddress()
	return peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.INBOUND),
		peer.WithCryptography(cryptography.CreateCommonCrypto()))
}

func fetchAvailablePeers() []peer.Peer {
	var remotePeers map[string]peer.Peer = make(map[string]peer.Peer)
	var currentPeers []peer.Peer = peer.GetPeers(string(""), string(""), string(""))
	for _, currentPeer := range currentPeers {
		if currentPeer.Role == peer.OUTBOUND {
			for _, peer := range synchronizeWithPeer(currentPeer) {
				remotePeers[peer.Hostname] = peer
			}
		}
	}
	var uniquePeers []peer.Peer = []peer.Peer{}
	for _, value := range remotePeers {
		uniquePeers = append(uniquePeers, value)
	}
	return uniquePeers
}

func synchronizeWithPeer(currentPeer peer.Peer) []peer.Peer {
	var synchronizedPeers []peer.Peer = []peer.Peer{}
	res := http.GET(
		currentPeer,
		currentPeer.Address+"/peer",
		&synchronizedPeers,
		"hostname",
		currentPeer.Hostname)
	if res == nil || res.StatusCode != http.OK {
		fmt.Printf("Peer %s is offline.\n", currentPeer.Hostname)
		return []peer.Peer{}
	}
	return synchronizedPeers
}
