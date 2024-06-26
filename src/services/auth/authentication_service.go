// Copyright (c) 2024 Berk Kirtay

package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Authenticate(authenticationModel AuthenticationModel, c *gin.Context) AuthenticationModel {
	if authenticationModel.Id == "" || authenticationModel.Password == "" {
		return authenticationModel
	}
	session := sessions.Default(c)
	var token string = generateToken(authenticationModel, retrieveTokenNonce(session))
	session.Set(token, authenticationModel.Id)
	session.Save()
	return CreateAuthenticationModel(
		WithId(authenticationModel.Id),
		WithName(authenticationModel.Name),
		WithToken(token))
}

func generateToken(authenticationModel AuthenticationModel, nonce string) string {
	sha256 := sha256.New()
	sha256.Write([]byte(authenticationModel.Id))
	sha256.Write([]byte(authenticationModel.Name))
	sha256.Write([]byte(authenticationModel.Password))
	sha256.Write([]byte(strconv.FormatInt(rand.Int63(), 10)))
	sha256.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
	sha256.Write([]byte(nonce))
	sha256.Write(sha256.Sum(nil))
	return base64.StdEncoding.EncodeToString(sha256.Sum(nil))
}

func retrieveTokenNonce(session sessions.Session) string {
	var currentNonce interface{} = session.Get("TokenNonce")
	if currentNonce == nil {
		currentNonce = int64(1000000000)
	}
	currentNonce = currentNonce.(int64) + 1
	session.Set("TokenNonce", currentNonce)
	session.Save()
	return strconv.FormatInt(currentNonce.(int64), 10)
}
