// Copyright (c) 2024 Berk Kirtay

package middlewares

import (
	"crypto/rand"
	"main/services/auth"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitializeSession(app *gin.Engine) {
	var key []byte = make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	store := cookie.NewStore([]byte("session"), key)
	//TODO store.Options
	app.Use(sessions.Sessions("user-session", store))
}

func ValidateAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAuthenticated(c) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": "Connection is not authorized."})
		}
		c.Next()
	}

}

func isAuthenticated(c *gin.Context) bool {
	session := sessions.Default(c)
	authToken := c.Request.Header.Get("Authorization")
	userId := c.Request.Header.Get("Session")
	publicKey := c.Request.Header.Get("PublicKey")
	if publicKey != "" {
		// var token string =
		auth.InitializeSessionWithDiffieHellman(c, publicKey, userId)
		// return token == authToken
	}
	sessionId := session.Get(authToken)
	return sessionId == userId
}
