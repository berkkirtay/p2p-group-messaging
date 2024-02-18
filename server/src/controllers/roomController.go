package controllers

import (
	"encoding/json"
	"main/services/room"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getRooms(c *gin.Context) {
	res := room.GetRooms(c.Query("id"), c.Query("size"))
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func createRoom(c *gin.Context) {
	var roomBody room.Room
	err := json.NewDecoder(c.Request.Body).Decode(&roomBody)
	if err != nil { //Error: Invalid body"
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
	}
	res := room.PostRoom(roomBody)
	c.JSON(http.StatusCreated, res)
}

func deleteRooms(c *gin.Context) {
	res := room.DeleteRooms(strings.Split(c.Query("id"), "&"))
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func updateRoom(c *gin.Context) {
	var roomBody room.Room
	err := json.NewDecoder(c.Request.Body).Decode(&roomBody)
	if err != nil { //Error: Invalid body"
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
	}
	res := room.UpdateRoom(c.Query("id"), roomBody)
	c.JSON(http.StatusOK, res)
}

func joinRoom(c *gin.Context) {
	var roomBody room.Room
	err := json.NewDecoder(c.Request.Body).Decode(&roomBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		return
	}
	res := room.JoinRoom(c.Query("id"), roomBody, c.Request.Header.Get("Session"))
	if res.Id == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": "You are not authorized to enter this room."})
	}
	c.JSON(http.StatusOK, res)
}

func leaveRoom(c *gin.Context) {
	res := room.LeaveRoom(c.Query("id"))
	c.JSON(http.StatusOK, res)
}

func ReceiveMessagesHTTP(c *gin.Context) {
	res := room.ReceiveMessages(c.Query("id"), c.Query("size"), c.Request.Header.Get("Session"))
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, res)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func SendAMessageHTTP(c *gin.Context) {
	var messageBody room.Message
	err := json.NewDecoder(c.Request.Body).Decode(&messageBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error:": err.Error()})
		return
	}
	res := room.SendAMessage(c.Query("id"), c.Request.Header.Get("Session"), messageBody)
	if res.Id == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error:": "Message could not be sent."})
	}
	c.JSON(http.StatusCreated, res)
}

func DeleteAMessageHTTP(c *gin.Context) {

}

func Roomouter(routerGroup *gin.RouterGroup) {
	roomRouter := routerGroup.Group("/room")
	{
		roomRouter.GET("", getRooms)
		roomRouter.GET("/messages", ReceiveMessagesHTTP)
		roomRouter.POST("", createRoom)
		roomRouter.POST("/messages", SendAMessageHTTP)
		roomRouter.DELETE("", deleteRooms)
		roomRouter.PUT("", updateRoom)
		roomRouter.POST("/join", joinRoom)
		roomRouter.POST("/leave", leaveRoom)
	}
}
