// Copyright (c) 2024 Berk Kirtay

package cryptography

type Encryption struct {
	key   string
	IV    int64
	nonce int64
}

func (encryption *Encryption) WithKey(key string) *Encryption {
	encryption.key = key
	return encryption
}

func (encryption *Encryption) WithIV(iV int64) *Encryption {
	encryption.IV = iV
	return encryption
}

func (encryption *Encryption) WithNonce(nonce int64) *Encryption {
	encryption.nonce = nonce
	return encryption
}

func (encryption *Encryption) Build() Encryption {
	return *encryption
}

func DefaultEncryption() *Encryption {
	return &Encryption{}
}
