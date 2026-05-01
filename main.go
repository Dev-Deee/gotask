package main

import (
	"gotask/config"
	"gotask/db"
	"gotask/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}
	cfg := config.Load()
	db.Connect(&cfg)

	router := gin.Default()
	routes.SetupRoutes(router)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(router.Run(":" + port))
	//router.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	//})
}
