// Copyright (c) 2024 Berk Kirtay

package peer

import (
	"context"
	"main/infra/store"

	"go.mongodb.org/mongo-driver/bson"
)

var repository = store.NewRepo("peer")
var masterPeer Peer

func GetPeers() []Peer {
	return retrievePeers()
}

func PostPeer(peer Peer) Peer {
	builtPeer := CreatePeer(
		WithHostname(peer.Hostname),
		WithName(peer.Name),
		WithAddress(peer.Address),
		WithRole(peer.Role),
		WithCryptography(peer.Cryptography))
	filter := bson.D{
		{Key: "hostname", Value: builtPeer.Hostname},
		{Key: "role", Value: builtPeer.Role}}
	cur, _ := repository.FindOne(filter, nil)
	if cur == nil {
		repository.InsertOne(builtPeer)
	}
	if peer.Role == INBOUND {
		masterPeer = peer
	}
	return builtPeer
}

func DeletePeer(hostname string) int64 {
	filter := bson.D{{Key: "hostname", Value: hostname}}
	res, _ := repository.DeleteMany(filter, nil)
	return res.DeletedCount
}

func retrievePeers() []Peer {
	var peers []Peer = []Peer{}
	list, _ := repository.Find(bson.D{{}}, nil)
	for list.Next(context.TODO()) {
		var currentPeer Peer
		err := list.Decode(&currentPeer)
		if err != nil {
			panic(err)
		}
		peers = append(peers, currentPeer)
	}
	return peers
}

func GetMasterPeer() Peer {
	return masterPeer
}
