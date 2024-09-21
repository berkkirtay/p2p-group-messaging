// Copyright (c) 2024 Berk Kirtay

package cryptography

type Cryptography struct {
	Sign       string    `json:"sign,omitempty" bson:"sign,omitempty"`
	PublicKey  string    `json:"publicKey,omitempty" bson:"publicKey,omitempty"`
	PrivateKey string    `json:"privateKey,omitempty" bson:"privateKey,omitempty"`
	Nonce      int64     `json:"nonce,omitempty" bson:"nonce,omitempty"`
	Timestamp  string    `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Hash       string    `json:"hash,omitempty" bson:"hash,omitempty"`
	Elliptic   *Elliptic `json:"elliptic,omitempty" bson:"elliptic,omitempty"`
}

type CryptographyOption func(Cryptography) Cryptography

func WithSign(sign string) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.Sign = sign
		return Cryptography
	}
}

func WithPublicKey(publicKey string) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.PublicKey = publicKey
		return Cryptography
	}
}

func WithPrivateKey(privateKey string) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.PrivateKey = privateKey
		return Cryptography
	}
}

func WithNonce(nonce int64) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.Nonce = nonce
		return Cryptography
	}
}

func WithTimestamp(timestamp string) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.Timestamp = timestamp
		return Cryptography
	}
}

func WithHash(hash string) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.Hash = hash
		return Cryptography
	}
}

func WithElliptic(elliptic *Elliptic) CryptographyOption {
	return func(Cryptography Cryptography) Cryptography {
		Cryptography.Elliptic = elliptic
		return Cryptography
	}
}

func CreateDefaultCryptography() Cryptography {
	return Cryptography{}
}

func CreateCryptography(options ...CryptographyOption) *Cryptography {
	Cryptography := CreateDefaultCryptography()

	for _, option := range options {
		Cryptography = option(Cryptography)
	}

	return &Cryptography
}
