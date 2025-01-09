// Copyright (c) 2024 Berk Kirtay

package peer

import (
	"context"
	"main/infra/store"
	"main/services/auth"

	"go.mongodb.org/mongo-driver/bson"
)

var repository = store.NewRepo("peer")

func GetPeers(hostname string, role string, userId string) []Peer {
	var peers []Peer = []Peer{}
	if hostname != "" && role != "" {
		var peer Peer = Peer{}
		filter := bson.D{{
			Key: "hostname", Value: hostname},
			{Key: "role", Value: role}}
		cur, err := repository.FindOne(filter, nil)
		if cur != nil && err == nil {
			cur.Decode(&peer)
			enrichWithPeerEllipticKeys(&peer, userId)
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
			enrichWithPeerEllipticKeys(&currentPeer, userId)
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
	return builtPeer
}

func DeletePeer(hostname string) int64 {
	filter := bson.D{{Key: "hostname", Value: hostname}}
	res, _ := repository.DeleteMany(filter, nil)
	return res.DeletedCount
}

func enrichWithPeerEllipticKeys(peer *Peer, userId string) {
	if userId != "" {
		elliptic := auth.GetDiffieHellmanUserAuthentication(userId)
		peer.Cryptography.Elliptic = &elliptic
	}
}
