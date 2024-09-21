// Copyright (c) 2024 Berk Kirtay

package auth

import (
	"main/infra/cryptography"
	"main/services/peer"
	"main/services/user"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authenticate(authBody AuthenticationModel, c *gin.Context) AuthenticationModel {
	var receivedUser user.User = user.CreateUser(
		user.WithId(authBody.Id),
		user.WithName(authBody.Name),
		user.WithPassword(authBody.Password),
		user.WithCryptography(authBody.Cryptography))
	var actualUser user.User = user.GetUser(receivedUser.Id, receivedUser.Name)
	if actualUser.IsPeer {
		return authenticateWithDiffieHellman(receivedUser, actualUser, c)
	} else {
		return authenticateWithPassword(receivedUser, actualUser, c)
	}
}

func authenticateWithPKCS(
	receivedUser user.User,
	actualUser user.User,
	c *gin.Context) AuthenticationModel {
	var verification bool = cryptography.VerifySignature(
		[]string{
			receivedUser.Name},
		receivedUser.Cryptography.Sign,
		receivedUser.Cryptography.PublicKey)
	if verification {
		var token string = initializeSessionForUser(c, actualUser)
		var encryptedToken string = cryptography.EncryptRSA(
			token,
			receivedUser.Cryptography.PublicKey)
		return CreateAuthenticationModel(
			WithId(actualUser.Id),
			WithName(actualUser.Name),
			WithToken(encryptedToken),
			WithCryptography(actualUser.Cryptography),
		)
	} else {
		return CreateDefaultAuthenticationModel()
	}
}

func authenticateWithDiffieHellman(
	receivedUser user.User,
	actualUser user.User,
	c *gin.Context) AuthenticationModel {
	var verification bool = cryptography.VerifySignature(
		[]string{
			receivedUser.Name},
		receivedUser.Cryptography.Sign,
		receivedUser.Cryptography.PublicKey)
	if verification && receivedUser.Name == actualUser.Name &&
		receivedUser.Cryptography.Sign == actualUser.Cryptography.Sign {
		InitializeSessionWithDiffieHellman(
			c,
			receivedUser.Cryptography.Elliptic.PublicKey,
			receivedUser.Id)
		crypto := cryptography.CreateCryptography(
			cryptography.WithPublicKey(peer.GetMasterPeer().Cryptography.PublicKey))
		return CreateAuthenticationModel(
			WithId(actualUser.Id),
			WithName(actualUser.Name),
			WithCryptography(crypto),
		)
	} else {
		return CreateDefaultAuthenticationModel()
	}
}

func authenticateWithPassword(
	receivedUser user.User,
	actualUser user.User,
	c *gin.Context) AuthenticationModel {
	if receivedUser.Id == "" || receivedUser.Password == "" {
		return CreateDefaultAuthenticationModel()
	}

	if actualUser.Id == "" {
		return CreateDefaultAuthenticationModel()
	}
	if actualUser.Password != receivedUser.Password {
		return CreateDefaultAuthenticationModel()
	}
	var token string = initializeSessionForUser(c, actualUser)
	return CreateAuthenticationModel(
		WithId(actualUser.Id),
		WithName(actualUser.Name),
		WithToken(token),
		WithCryptography(actualUser.Cryptography),
	)
}

func initializeSessionForUser(c *gin.Context, user user.User) string {
	session := sessions.Default(c)
	var token string = generateToken(user, getTokenNonce(session))
	session.Set(token, user.Id)
	session.Save()
	return token
}

func InitializeSessionWithDiffieHellman(c *gin.Context, publicKey string, userId string) {
	session := sessions.Default(c)
	var token string = cryptography.DiffieHellman(
		peer.GetMasterPeer().Cryptography.Elliptic.PrivateKey, publicKey)
	session.Set(token, userId)
	session.Save()
}

func generateToken(user user.User, nonce string) string {
	return cryptography.GenerateEncodedSHA256([]string{
		user.Id,
		user.Name,
		user.Password,
		strconv.FormatInt(rand.Int63(), 10),
		strconv.FormatInt(time.Now().Unix(), 10),
		nonce})
}

func getTokenNonce(session sessions.Session) string {
	var currentNonce interface{} = session.Get("TokenNonce")
	if currentNonce == nil {
		currentNonce = cryptography.GenerateANonce()
	}
	currentNonce = currentNonce.(int64) + 1
	session.Set("TokenNonce", currentNonce)
	session.Save()
	return strconv.FormatInt(currentNonce.(int64), 10)
}
