package routes

import (
	"github.com/gin-gonic/gin"

	"gotask/handlers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", handlers.GetTasks)
			tasks.GET("/:id", handlers.GetTask)
			tasks.POST("", handlers.CreateTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}
	}
}
