// Copyright (c) 2024 Berk Kirtay

package controllers

import (
	"encoding/json"
	"main/services/peer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPeers(c *gin.Context) {
	res := peer.GetPeers(c.Query("hostname"), c.Query("role"))
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
	if res.Hostname == "" {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusCreated, res)
	}
}

func deletePeer(c *gin.Context) {
	deletedCount := peer.DeletePeer(c.Query("hostId"))
	if deletedCount == 0 {
		c.Status(http.StatusNotFound)
	} else {
		c.Status(http.StatusOK)
	}
}

func PeerRouter(routerGroup *gin.RouterGroup) {
	peerRouter := routerGroup.Group("/peer")
	{
		peerRouter.GET("", getPeers)
		peerRouter.POST("", postPeer)
		peerRouter.DELETE("", deletePeer)
	}
}
