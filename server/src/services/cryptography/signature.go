package cryptography

type Signature struct {
	Sign      string `json:"sign,omitempty" bson:"sign,omitempty"`
	PublicKey string `json:"publicKey,omitempty" bson:"publicKey,omitempty"`
	Nonce     int64  `json:"nonce,omitempty" bson:"nonce,omitempty"`
	Timestamp string `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	Hash      string `json:"hash,omitempty" bson:"hash,omitempty"`
}

type SignatureOption func(Signature) Signature

func WithSign(sign string) SignatureOption {
	return func(signature Signature) Signature {
		signature.Sign = sign
		return signature
	}
}

func WithPublicKey(publicKey string) SignatureOption {
	return func(signature Signature) Signature {
		signature.PublicKey = publicKey
		return signature
	}
}

func WithNonce(nonce int64) SignatureOption {
	return func(signature Signature) Signature {
		signature.Nonce = nonce
		return signature
	}
}

func WithTimestamp(timestamp string) SignatureOption {
	return func(signature Signature) Signature {
		signature.Timestamp = timestamp
		return signature
	}
}

func WithHash(hash string) SignatureOption {
	return func(signature Signature) Signature {
		signature.Hash = hash
		return signature
	}
}

func CreateDefaultSignature() Signature {
	return Signature{}
}

func CreateSignature(options ...SignatureOption) *Signature {
	signature := CreateDefaultSignature()

	for _, option := range options {
		signature = option(signature)
	}

	return &signature
}
