// Copyright (c) 2024 Berk Kirtay

package controllers

import (
	"main/services/info"
	"net/http"

	"github.com/gin-gonic/gin"
)

func postInfo(c *gin.Context) {
	res := info.HandlePeerInfo(c.Query("type"))
	if res {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNoContent)
	}
}

func InfoRouter(routerGroup *gin.RouterGroup) {
	peerRouter := routerGroup.Group("/info")
	{
		peerRouter.POST("", postInfo)
	}
}
