// Copyright (c) 2024 Berk Kirtay

package controllers

import (
	"encoding/json"
	"main/services/auth"
	"main/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postAuthRequest(c *gin.Context) {
	var userBody user.User
	err := json.NewDecoder(c.Request.Body).Decode(&userBody)
	if err != nil {
		panic(err)
	}
	res := auth.Authenticate(userBody, c)
	if res.Token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	c.JSON(http.StatusAccepted, res)
}

func AuthRouter(routerGroup *gin.RouterGroup) {
	authRouter := routerGroup.Group("/auth")
	{
		authRouter.POST("", postAuthRequest)
	}
}
