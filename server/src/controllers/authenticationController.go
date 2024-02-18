package controllers

import (
	"encoding/json"
	"main/auth"
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

	res := auth.Authenticate(auth.CreateAuthenticationModel(
		auth.WithId(userBody.Id),
		auth.WithPassword(userBody.Password)), c)
	if res.Token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
	}
	c.JSON(http.StatusAccepted, res)
}

func AuthRouter(routerGroup *gin.RouterGroup) {
	authRouter := routerGroup.Group("/auth")
	{
		authRouter.POST("", postAuthRequest)
	}
}
