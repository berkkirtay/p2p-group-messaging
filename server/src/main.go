package main

import (
	"main/middlewares"

	"github.com/gin-gonic/gin"
)

const ADDRESS string = ":8080"
const API string = "/api/"

func main() {
	app := gin.New()
	app.Use(gin.Logger())
	middlewares.InitializeSession(app)
	router := app.Group(API)
	middlewares.InitializeRouters(router)
	app.Run(ADDRESS)
}
