package handlers

import (
	"gotask/db"
	"gotask/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task

	user, _ := c.MustGet("currentUser").(models.User)
	query := db.DB.Where("user_id = ?", user.ID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func CreateTask(c *gin.Context) {
	var input models.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.MustGet("currentUser").(models.User)

	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		UserID:      user.ID,
	}
	if err := db.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": task})
}

func GetTask(c *gin.Context) {
	user, _ := c.MustGet("currentUser").(models.User)
	var task models.Task
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func UpdateTask(c *gin.Context) {
	user, _ := c.MustGet("currentUser").(models.User)

	var task models.Task
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&task).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": task})
}

func GetTasksByStatus(c *gin.Context) {
	user, _ := c.MustGet("currentUser").(models.User)

	var tasks []models.Task
	status := c.Query("status")

	query := db.DB.Where("user_id = ?", user.ID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
	})
}

func DeleteTask(c *gin.Context) {
	user, _ := c.MustGet("currentUser").(models.User)

	var task models.Task
	if err := db.DB.Where("id = ? AND user_id = ?", c.Param("id"), user.ID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	if task.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "not your task"})
		return
	}
	db.DB.Delete(&task)
	c.JSON(http.StatusNoContent, nil)
}

func PatchTask(c *gin.Context) {
	user, _ := c.MustGet("currentUser").(models.User)

	var task models.Task

	if err := db.DB.
		Where("id = ? AND user_id = ?", c.Param("id"), user.ID).
		First(&task).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "task not found",
		})
		return
	}

	var input models.PatchTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&task).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}
