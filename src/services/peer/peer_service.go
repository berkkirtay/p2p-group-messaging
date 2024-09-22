// Copyright (c) 2024 Berk Kirtay

package peer

import (
	"context"
	"main/infra/cryptography"
	"main/infra/store"

	"go.mongodb.org/mongo-driver/bson"
)

var repository = store.NewRepo("peer")
var masterPeer Peer

func GetPeers(hostname string, role string) []Peer {
	var peers []Peer = []Peer{}
	if hostname != "" && role != "" {
		var peer Peer = Peer{}
		filter := bson.D{{
			Key: "hostname", Value: hostname},
			{Key: "role", Value: role}}
		cur, err := repository.FindOne(filter, nil)
		if cur != nil && err == nil {
			cur.Decode(&peer)
			UpdatePeerEllipticKeys(&peer)
			peers = append(peers, peer)
		}
	} else {
		list, _ := repository.Find(bson.D{{}}, nil)
		for list.Next(context.TODO()) {
			var currentPeer Peer
			err := list.Decode(&currentPeer)
			if err != nil {
				panic(err)
			}
			UpdatePeerEllipticKeys(&currentPeer)
			peers = append(peers, currentPeer)
		}
	}
	return peers
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

func GetMasterPeer() Peer {
	return masterPeer
}

func UpdatePeerEllipticKeys(peer *Peer) {
	if peer.Cryptography.Elliptic.PrivateKey == nil {
		peer.Cryptography.Elliptic = cryptography.CreateElliptic(
			cryptography.WithEllipticKeys(cryptography.GenerateEllipticCurveKeys()))
	}
	if peer.Role == INBOUND {
		masterPeer = *peer
	}
}
