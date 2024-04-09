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
	var foundUser = user.GetUser(userBody.Id, userBody.Name)
	if foundUser.Id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if foundUser.Password != userBody.Password {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	res := auth.Authenticate(auth.CreateAuthenticationModel(
		auth.WithId(foundUser.Id),
		auth.WithName(foundUser.Name),
		auth.WithPassword(foundUser.Password)), c)
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
