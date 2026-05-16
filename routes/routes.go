package routes

import (
	"gotask/middleware"

	"github.com/gin-gonic/gin"

	"gotask/handlers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.RegisterUser)
		auth.POST("/login", handlers.LoginUser)
	}
	tasks := api.Group("/tasks")
	tasks.Use(middleware.RequireAuth)
	{
		tasks.GET("", handlers.GetTasks)
		tasks.GET("/:id", handlers.GetTask)
		tasks.POST("", handlers.CreateTask)
		tasks.PUT("/:id", handlers.UpdateTask)
		tasks.DELETE("/:id", handlers.DeleteTask)
	}
}
