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

var thisPeer peer.Peer
var assignedPeer peer.Peer

func InitializePeer(identifier string) {
	network.InitializeBroadcast(identifier)
	network.StartPeerBroadcast()
	handleUserPeerSelection()
	finalizeInitialization()
	go peerListenerLoop()
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

func GetMainPeer() peer.Peer {
	return assignedPeer
}

func GetNode() peer.Peer {
	return thisPeer
}

func peerListenerLoop() {
	for {
		targetHostname, targetAddress := network.ListenForPeerBroadcast()
		registerPeer(targetAddress, targetHostname)
	}
}

func registerPeer(targetAddress string, targetHostname string) {
	var newPeer peer.Peer = peer.CreatePeer(
		peer.WithHostname(assignedPeer.Hostname),
		peer.WithAddress(assignedPeer.Address),
		peer.WithRole(peer.OUTBOUND),
		peer.WithCryptography(assignedPeer.Cryptography))
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
		peer.WithHostname(targetHostname),
		peer.WithAddress(targetAddress),
		peer.WithRole(peer.INBOUND)))
}

func handleUserPeerSelection() {
	var peers []peer.Peer = fetchOutboundPeers()
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
			assignedPeer = peers[number]
			return
		}
	}
	fmt.Println("No connection has been done to any peer. Initializing this node as an active peer.")
	assignedPeer = initializeAMasterPeer()
}

func initializeAMasterPeer() peer.Peer {
	hostname, address := network.GetHostAddress()
	return peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.OUTBOUND),
		peer.WithCryptography(cryptography.CreateCommonCrypto()))
}

func finalizeInitialization() {
	if thisPeer.Address == assignedPeer.Address {
		thisPeer = peer.CreatePeer(peer.WithPeer(assignedPeer))
	} else {
		hostname, address := network.GetHostAddress()
		thisPeer = peer.CreatePeer(
			peer.WithHostname(hostname),
			peer.WithAddress(address))
	}
	peer.PostPeer(assignedPeer)
	fmt.Printf("A master peer is initialized: %s\n", assignedPeer.Hostname)
}

func fetchOutboundPeers() []peer.Peer {
	return fetchAvailablePeers(peer.OUTBOUND)
}

func fetchInboundPeers() []peer.Peer {
	return fetchAvailablePeers(peer.INBOUND)
}

func fetchAvailablePeers(peerRole string) []peer.Peer {
	var remotePeers map[string]peer.Peer = make(map[string]peer.Peer)
	var currentPeers []peer.Peer = peer.GetPeers(string(""), string(""), string(""))
	for _, currentPeer := range currentPeers {
		if currentPeer.Role == peerRole { //
			for _, peer := range synchronizeWithPeer(currentPeer) {
				if peer.Role == peerRole {
					remotePeers[peer.Address] = peer
				}
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
		peer.CreateDefaultPeer(),
		currentPeer.Address+"/peer",
		&synchronizedPeers,
		"hostname",
		currentPeer.Hostname)
	if res == nil || res.StatusCode != http.OK {
		//fmt.Printf("Peer %s is offline.\n", currentPeer.Hostname)
		return []peer.Peer{}
	}
	return synchronizedPeers
}
