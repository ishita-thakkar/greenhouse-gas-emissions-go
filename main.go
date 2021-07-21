package main

import (
	"bluesky.com/greenhouse-gas-emissions/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/ping", controller.PingHandler)
	router.GET("/countries", controller.Countries)
	router.GET("/country/:id", controller.Country)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
