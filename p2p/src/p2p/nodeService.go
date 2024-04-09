package p2p

import (
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

func StartNode() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	//app.Use(gin.Logger())
	middlewares.InitializeSession(app)
	router := app.Group(API)
	middlewares.InitializeRouters(router)
	app.Run(ADDRESS)
}
