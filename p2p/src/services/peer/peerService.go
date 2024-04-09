package peer

import (
	"context"
	"main/infrastructure"

	"go.mongodb.org/mongo-driver/bson"
)

var repository = infrastructure.NewRepo("peer")

func GetPeers() []Peer {
	return retrievePeers()
}

func PostPeer(peer Peer) []Peer {
	builtPeer := CreatePeer(
		WithHostId(peer.HostId),
		WithName(peer.Name),
		WithAddress(peer.Address),
		WithRole(peer.Role))
	repository.InsertOne(context.TODO(), builtPeer)
	return retrievePeers()
}

func retrievePeers() []Peer {
	var peers []Peer = []Peer{}
	list, _ := repository.Find(context.TODO(), bson.D{{}}, nil)
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
