package controllers

import "github.com/gin-gonic/gin"

func postAuthorization(c *gin.Context) {
	c.SetAccepted()
}

func AuthorizationRouter(routerGroup *gin.RouterGroup) {
	authorizationRouter := routerGroup.Group("/authorization")
	{
		authorizationRouter.POST("", postAuthorization)
	}
}
