// Copyright (c) 2024 Berk Kirtay

package peer

type Peer struct {
	HostId  string `json:"hostId,omitempty" bson:"hostId,,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Address string `json:"address,omitempty" bson:"address,omitempty"`
	Role    string `json:"role,omitempty" bson:"role,omitempty"`
}

type PeerOption func(Peer) Peer

func WithHostId(hostId string) PeerOption {
	return func(peer Peer) Peer {
		peer.HostId = hostId
		return peer
	}
}

func WithName(name string) PeerOption {
	return func(peer Peer) Peer {
		peer.Name = name
		return peer
	}
}

func WithAddress(address string) PeerOption {
	return func(peer Peer) Peer {
		peer.Address = address
		return peer
	}
}

func WithRole(role string) PeerOption {
	return func(peer Peer) Peer {
		peer.Role = role
		return peer
	}
}

func CreateDefaultPeer() Peer {
	return Peer{}
}

func CreatePeer(options ...PeerOption) Peer {
	peer := CreateDefaultPeer()

	for _, option := range options {
		peer = option(peer)
	}

	return peer
}
