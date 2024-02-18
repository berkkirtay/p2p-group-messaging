package controllers

import (
	"github.com/gin-gonic/gin"
)

func postCommand(c *gin.Context) {
	c.SetAccepted()
}

func CommandRouter(routerGroup *gin.RouterGroup) {
	commandRouter := routerGroup.Group("/command")
	{
		commandRouter.POST("/blockchain", postCommand)
	}
}
