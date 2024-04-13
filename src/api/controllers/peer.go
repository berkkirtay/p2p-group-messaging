// Copyright (c) 2024 Berk Kirtay

package controllers

import (
	"encoding/json"
	"main/services/peer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPeers(c *gin.Context) {
	res := peer.GetPeers()
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func postPeer(c *gin.Context) {
	var peerBody peer.Peer
	err := json.NewDecoder(c.Request.Body).Decode(&peerBody)
	if err != nil {
		panic(err)
	}
	res := peer.PostPeer(peerBody)
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func PeerRouter(routerGroup *gin.RouterGroup) {
	peerRouter := routerGroup.Group("/peer")
	{
		peerRouter.GET("", getPeers)
		peerRouter.POST("", postPeer)
	}
}
