// Copyright (c) 2024 Berk Kirtay

package controllers

import (
	"encoding/json"
	"main/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postAuthRequest(c *gin.Context) {
	var authBody auth.AuthenticationModel
	err := json.NewDecoder(c.Request.Body).Decode(&authBody)
	if err != nil {
		panic(err)
	}
	res := auth.Authenticate(authBody, c)
	if res.Id == "" {
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
