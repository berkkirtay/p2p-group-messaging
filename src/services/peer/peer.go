// Copyright (c) 2024 Berk Kirtay

package peer

const (
	INBOUND  = "INBOUND"
	OUTBOUND = "OUTBOUND"
)

type Peer struct {
	Hostname string `json:"hostname,omitempty" bson:"hostname,,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Address  string `json:"address,omitempty" bson:"address,omitempty"`
	Role     string `json:"role,omitempty" bson:"role,omitempty"`
}

type PeerOption func(Peer) Peer

func WithHostname(hostname string) PeerOption {
	return func(peer Peer) Peer {
		peer.Hostname = hostname
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

func WithPeer(newPeer Peer) PeerOption {
	return func(peer Peer) Peer {
		peer.Hostname = newPeer.Hostname
		peer.Name = newPeer.Name
		peer.Address = newPeer.Address
		peer.Role = newPeer.Role
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
