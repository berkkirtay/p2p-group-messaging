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
	if len(remotePeers) > 0 {
		handleUserPeerSelection(remotePeers)
	}

	if assignedPeer.Address == "" {
		fmt.Println("No remote peer is available. Making yourself an active peer.")
		assignedPeer = initializeAMasterPeer()
	}
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

func initializeAMasterPeer() peer.Peer {
	hostname, address := network.GetHostAddress()
	return peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.INBOUND),
		peer.WithCryptography(cryptography.CreateCommonCrypto()))
}

// TODO review here:
func registerPeer(targetAddress string, hostname string, address string) {
	var newPeer peer.Peer = peer.CreatePeer(
		peer.WithHostname(hostname),
		peer.WithAddress(address),
		peer.WithRole(peer.OUTBOUND),
		peer.WithCryptography(cryptography.CreateCommonCrypto()))
	// ,
	// peer.WithCryptography(assignedPeer.Cryptography)
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
		peer.WithRole(peer.INBOUND)))

	// ,
	// peer.WithCryptography(
	// 	cryptography.CreateCryptography(
	// 		cryptography.WithElliptic(assignedPeer.Cryptography.Elliptic)))
}

func handleUserPeerSelection(peers []peer.Peer) {
	fmt.Println("--------------")
	fmt.Println("i |	Hostname	|	Address		")
	for i, peer := range peers {
		fmt.Printf("%d |	%s	|	%s		\n", i, peer.Hostname, peer.Address)
	}
	fmt.Println("--------------")

	fmt.Println("Please select a peer to make a connection: ")
	var number int
	_, err := fmt.Scanf("%d", &number)
	if err != nil {
		panic(err)
	}
	if number <= len(peers)-1 {
		assignedPeer = peers[number]
	}
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
