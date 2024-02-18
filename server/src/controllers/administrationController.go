package controllers

import "github.com/gin-gonic/gin"

func getAdmins(c *gin.Context) {
	c.SetAccepted()
}

func AdministrationRouter(routerGroup *gin.RouterGroup) {
	administrationRouter := routerGroup.Group("/admin")
	{
		administrationRouter.GET("", getAdmins)
	}
}
