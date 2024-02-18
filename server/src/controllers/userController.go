package controllers

import (
	"encoding/json"
	"fmt"
	"main/services/audit"
	"main/services/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) {
	res := user.GetUsers(c.Query("id"), c.Query("size"))
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func postUser(c *gin.Context) {
	var userBody user.User
	err := json.NewDecoder(c.Request.Body).Decode(&userBody)
	if err != nil { //Error: Invalid body"
		panic(err)
	}
	userBody.Audit = audit.CreateAuditForUser(c.ClientIP())
	res := user.PostUser(userBody)
	fmt.Println(userBody.Audit)
	c.JSON(http.StatusCreated, res)
}

func UserRouter(routerGroup *gin.RouterGroup) {
	userRouter := routerGroup.Group("/users")
	{
		userRouter.GET("", getUser)
		userRouter.POST("", postUser)
	}
}
