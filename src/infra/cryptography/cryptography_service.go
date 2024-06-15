// Copyright (c) 2024 Berk Kirtay

package cryptography

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"io"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/zenazn/pkcs7pad"
)

const p float64 = 23
const g float64 = 5

var publicKey string
var privateKey string

func InitializeService() {
	// rand.Prime()
}

// TODO area.
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

func GenerateARandomMasterSecret() string {
	sha256 := crypto.SHA256.New()
	randValue, err := rand.Int(rand.Reader, big.NewInt(32))
	if err != nil {
		panic(err)
	}

	sha256.Write([]byte(strconv.FormatInt(randValue.Int64(), 10)))
	randomSecret := sha256.Sum([]byte("secret"))[0:16]

	encodedSecret := b64.StdEncoding.EncodeToString([]byte(randomSecret))
	return encodedSecret
}

func ServerSideDiffieHelmanKeyExhange(
	clientSecret string) []string {
	computedSecretOfClient, err := strconv.Atoi(clientSecret)
	if err != nil {
		panic(err)
	}
	computedSecretOfServer, err := strconv.Atoi(privateKey)
	if err != nil {
		panic(err)
	}
	commonSecretKey := math.Mod(math.Pow(float64(computedSecretOfClient), float64(computedSecretOfServer)), p)
	serverSecret := math.Mod(math.Pow(g, float64(computedSecretOfServer)), p)
	return []string{strconv.FormatFloat(commonSecretKey, 'f', -1, 64), strconv.FormatFloat(serverSecret, 'f', -1, 64)}
}

func PeerToPeerDiffieHelmanKeyExhange(p float64, g float64, privateKeys ...string) []string {
	combinedKeys := []string{}
	columns := []float64{}
	for _, privateKey := range privateKeys {
		secret, err := strconv.Atoi(privateKey)
		if err != nil {
			panic(err)
		}
		columns = append(columns, math.Mod(math.Pow(g, float64(secret)), p))
	}
	return combinedKeys
}

/*
 * Encryption functions take a string (preferably as base64) as input
 *  and returns an encrypted/decrypted string conventionally.
 */
func Encrypt(message string, key string) string {
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

func Decrypt(cipherText string, key string) string {
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
	if err != nil {
		panic(err)
	}
	return string(plainText[aes.BlockSize:])
}
