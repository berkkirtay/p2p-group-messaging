package cryptography

import "time"

func CreateDefaultCrypto(keyPair string, values ...interface{}) *Signature {
	hash, nonce := generateHash(values)
	return CreateSignature(WithPublicKey(keyPair),
		WithHash(hash),
		WithNonce(nonce),
		WithSign(generateSignature(keyPair, hash)),
		WithTimestamp(time.Now().Format(time.RFC1123)))
}

func generateHash(values []interface{}) (string, int64) {
	return values[0].(string), 1
}

func generateSignature(keyPair string, hash string) string {
	return hash
}
