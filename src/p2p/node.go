// Copyright (c) 2024 Berk Kirtay

package p2p

import (
	"fmt"
	"main/api/middlewares"
	"main/commands"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ADDRESS       string = ":"
	DEFAULT_PORT  string = "8080"
	RETRY_ADDRESS string = ":8081"
	API           string = "/api"
)

var nodeIdentifier string = ""

func StartClient() {
	commands.InitializePeer(nodeIdentifier)
	commands.InitializeCommandLine()
	commands.HandleInput()
}

func StartNode() {
	nodeIdentifier = ADDRESS
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	//	app.Use(gin.Logger())
	middlewares.InitializeSession(app)
	router := app.Group(API)
	middlewares.InitializeRouters(router)
	startWithIDentifier(app)
}

func startWithIDentifier(app *gin.Engine) {
	wait := make(chan bool)
	go func() {
		var nextPort, _ = strconv.Atoi(DEFAULT_PORT)
		for {
			nodeIdentifier = ADDRESS + strconv.Itoa(nextPort)
			err := app.Run(nodeIdentifier)
			if err == nil {
				wait <- true
				return
			}
			nextPort++
			fmt.Printf("DEFAULT PORT is already being used. Trying the next port: %d\n", nextPort)
		}

	}()
	go func() {
		time.Sleep(2000 * time.Millisecond)
		wait <- true
	}()
	<-wait
}
