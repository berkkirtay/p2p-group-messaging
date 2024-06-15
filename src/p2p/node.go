// Copyright (c) 2024 Berk Kirtay

package p2p

import (
	"main/api/middlewares"

	"github.com/gin-gonic/gin"
)

func StartNode() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	//app.Use(gin.Logger())
	middlewares.InitializeSession(app)
	router := app.Group(API)
	middlewares.InitializeRouters(router)
	go app.Run(ADDRESS)
}
