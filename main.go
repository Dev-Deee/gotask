package main

import (
	"gotask/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	routes.SetupRoutes(router)
	//router.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	//})
	router.Run(":8080")
}
