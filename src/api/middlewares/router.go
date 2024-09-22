// Copyright (c) 2024 Berk Kirtay

package middlewares

import (
	"main/api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRouters(routerGroup *gin.RouterGroup) {
	routerGroup.Use(gin.CustomRecovery(handleGenericPanic))
	controllers.PeerRouter(routerGroup)
	controllers.UserRouter(routerGroup)
	controllers.AuthRouter(routerGroup)
	routerGroup.Use(ValidateAuthentication())
	controllers.Roomouter(routerGroup)
}

func handleGenericPanic(c *gin.Context, err any) {
	defer func() {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{"Error:": err})
	}()
}
