// Copyright (c) 2024 Berk Kirtay

package cryptography

import "crypto/ecdh"

type Elliptic struct {
	PublicKey  string `json:"publicKey,omitempty" bson:"publicKey,omitempty"`
	PrivateKey *ecdh.PrivateKey
}

type EllipticOption func(Elliptic) Elliptic

func WithEllipticKeys(privateKey *ecdh.PrivateKey, publicKey string) EllipticOption {
	return func(Elliptic Elliptic) Elliptic {
		Elliptic.PrivateKey = privateKey
		Elliptic.PublicKey = publicKey
		return Elliptic
	}
}

func WithEllipticPublicKey(publicKey string) EllipticOption {
	return func(Elliptic Elliptic) Elliptic {
		Elliptic.PublicKey = publicKey
		return Elliptic
	}
}

func WithEllipticPrivateKey(privateKey *ecdh.PrivateKey) EllipticOption {
	return func(Elliptic Elliptic) Elliptic {
		Elliptic.PrivateKey = privateKey
		return Elliptic
	}
}

func CreateDefaultElliptic() Elliptic {
	return Elliptic{}
}

func CreateElliptic(options ...EllipticOption) *Elliptic {
	Elliptic := CreateDefaultElliptic()

	for _, option := range options {
		Elliptic = option(Elliptic)
	}

	return &Elliptic
}
