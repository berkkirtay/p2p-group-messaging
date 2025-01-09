// Copyright (c) 2024 Berk Kirtay

package cryptography

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"io"
	"math/big"
	"strconv"
	"time"

	"github.com/zenazn/pkcs7pad"
)

const (
	RSA_KEY_SIZE = 2048
)

func CreateCommonCrypto(values ...string) *Cryptography {
	hash := GenerateEncodedSHA256([]string(values))
	privateKey, publicKey := generateKeyPair()
	crypto := CreateCryptography(
		WithPublicKey(publicKey),
		WithPrivateKey(privateKey),
		WithHash(hash),
		WithNonce(GenerateANonce()),
		WithSign(generateSignature(privateKey, hash)),
		WithTimestamp(time.Now().Format(time.RFC1123)),
		WithElliptic(CreateElliptic(WithEllipticKeys(GenerateEllipticCurveKeys()))))
	return crypto
}

func generateSignature(privateKey string, hash string) string {
	decodedKey, err := b64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	decodedBlock, _ := pem.Decode(decodedKey)
	key, err := x509.ParsePKCS1PrivateKey(decodedBlock.Bytes)
	if err != nil {
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		key,
		crypto.SHA256,
		decodeSHA256(hash),
	)
	if err != nil {
		panic(err)
	}
	return b64.StdEncoding.EncodeToString(signature)
}

func VerifySignature(data []string, signature string, publicKey string) bool {
	decodedSignature, err := b64.StdEncoding.DecodeString(signature)
	if err != nil {
		panic(err)
	}
	decodedKey, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	decodedBlock, _ := pem.Decode(decodedKey)
	key, err := x509.ParsePKCS1PublicKey(decodedBlock.Bytes)
	if err != nil {
		panic(err)
	}
	calculatedHash := generateSHA256Object(data)
	err = rsa.VerifyPKCS1v15(key, crypto.SHA256, calculatedHash[:], decodedSignature)
	return err == nil
}

func GenerateEncodedSHA256(values []string) string {
	hash := generateSHA256Object(values)
	encodedHash := b64.StdEncoding.EncodeToString(hash[:])
	return encodedHash
}

func decodeSHA256(hash string) []byte {
	decodedHash, err := b64.StdEncoding.DecodeString(hash)
	if err != nil {
		panic(err)
	}
	return decodedHash
}

func generateSHA256Object(values []string) [32]byte {
	var data []byte = []byte{}
	for _, value := range values {
		data = append(data, []byte(value)...)
	}
	return sha256.Sum256(data)
}

func GenerateARandomMasterSecret() string {
	randValue, err := rand.Int(rand.Reader, big.NewInt(32))
	if err != nil {
		panic(err)
	}
	return GenerateEncodedSHA256([]string{strconv.FormatInt(randValue.Int64(), 10)})
}

func EnrichMasterSecret(secret string, hash string) string {
	randValue, err := rand.Int(rand.Reader, big.NewInt(32))
	if err != nil {
		panic(err)
	}
	return GenerateEncodedSHA256(
		[]string{
			strconv.FormatInt(randValue.Int64(), 10),
			secret,
			hash},
	)
}

func generateKeyPair(salts ...string) (string, string) {
	keyPair, err := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)
	if err != nil {
		panic(err)
	}

	err = keyPair.Validate()
	if err != nil {
		panic(err)
	}

	encodedPrivateKey := b64.StdEncoding.EncodeToString(
		pem.EncodeToMemory(&pem.Block{Bytes: x509.MarshalPKCS1PrivateKey(keyPair)}))

	encodedPublicKey := b64.StdEncoding.EncodeToString(
		pem.EncodeToMemory(&pem.Block{Bytes: x509.MarshalPKCS1PublicKey(&keyPair.PublicKey)}))
	return encodedPrivateKey, encodedPublicKey
}

func EncryptRSA(data string, publicKey string) string {
	decodedKey, err := b64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	decodedBlock, _ := pem.Decode(decodedKey)

	key, err := x509.ParsePKCS1PublicKey(decodedBlock.Bytes)
	if err != nil {
		panic(err)
	}
	ciperText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, key, []byte(data), nil)
	if err != nil {
		panic(err)
	}
	return b64.StdEncoding.EncodeToString(ciperText)
}

func DecryptRSA(data string, privateKey string) string {
	decodedData, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		panic(err)
	}
	decodedKey, err := b64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}
	decodedBlock, _ := pem.Decode(decodedKey)
	key, err := x509.ParsePKCS1PrivateKey(decodedBlock.Bytes)
	if err != nil {
		panic(err)
	}
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, decodedData, nil)
	if err != nil {
		panic(err)
	}
	return string(plainText)
}

/*
 * Encryption functions take a string (preferably as base64) as input
 * and returns an encrypted/decrypted string conventionally.
 */
func EncryptAES(message string, key string) string {
	decodedKey, _ := b64.StdEncoding.DecodeString(key)
	encryptor, err := aes.NewCipher(decodedKey)
	if err != nil {
		panic(err)
	}
	var paddedMessage []byte = pkcs7pad.Pad([]byte(message), aes.BlockSize)
	ciphertext := make([]byte, len(paddedMessage)+aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		panic(err)
	}
	cbc := cipher.NewCBCEncrypter(encryptor, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], paddedMessage)
	ciphertext = append(ciphertext, iv...)
	return b64.StdEncoding.EncodeToString(ciphertext)
}

func DecryptAES(cipherText string, key string) string {
	decodedCipherText, _ := b64.StdEncoding.DecodeString(cipherText)
	decodedKey, _ := b64.StdEncoding.DecodeString(key)
	decryptor, err := aes.NewCipher(decodedKey)
	if err != nil {
		panic(err)
	}
	iv := decodedCipherText[:aes.BlockSize]
	cbc := cipher.NewCBCDecrypter(decryptor, iv)
	decryptedText := decodedCipherText[0 : len(decodedCipherText)-aes.BlockSize]
	plainText := make([]byte, len(decryptedText))
	cbc.CryptBlocks(plainText, decryptedText)
	plainText, err = pkcs7pad.Unpad(plainText)
	// todo plainText is empty so thats why we get unPadding error
	if err != nil {
		panic(err)
	}
	return string(plainText[aes.BlockSize:])
}

func GenerateEllipticCurveKeys() (*ecdh.PrivateKey, string) {
	key, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	publicKey, err := x509.MarshalPKIXPublicKey(key.PublicKey())
	if err != nil {
		panic(err)
	}
	return key, b64.StdEncoding.EncodeToString(publicKey)
}

func DiffieHellman(receiverKey *ecdh.PrivateKey, senderKey string) string {
	decodedPublicKey, err := b64.StdEncoding.DecodeString(senderKey)
	if err != nil {
		panic(err)
	}
	publicKey, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		panic(err)
	}
	sharedSecret, err := receiverKey.ECDH(publicKey.(*ecdh.PublicKey))
	if err != nil {
		panic(err)
	}
	return GenerateEncodedSHA256([]string{string(sharedSecret)})
}

// TODO
func GenerateANonce() int64 {
	randValue, err := rand.Int(rand.Reader, big.NewInt(64))
	if err != nil {
		panic(err)
	}
	return randValue.Int64()
}
