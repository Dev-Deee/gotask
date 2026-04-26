package handlers

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"gotask/models"

	"github.com/gin-gonic/gin"
)

var (
	tasks  []models.Task
	nextID uint64 = 1
	mu     sync.RWMutex
)

func GetTasks(c *gin.Context) {
	mu.Lock()

	defer mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"data": tasks})
}

func CreateTask(c *gin.Context) {
	var input models.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	newTask := models.Task{
		ID:          nextID,
		Title:       input.Title,
		Description: input.Description,
		Status:      models.StatusPending,
		CreatedAt:   time.Now(),
	}
	nextID++
	tasks = append(tasks, newTask)
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for _, task := range tasks {
		if task.ID == uint64(id) {
			c.JSON(http.StatusOK, gin.H{"data": task})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == uint64(id) {
			tasks[i].Title = input.Title
			tasks[i].Description = input.Description
			tasks[i].Status = input.Status

			c.JSON(http.StatusOK, gin.H{"data": tasks[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
	}
	mu.Lock()
	defer mu.Unlock()

	for i, task := range tasks {
		if task.ID == uint64(id) {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}

func GetTasksByStatus(c *gin.Context) {
	statusQuery := c.Query("status")

	mu.RLock()
	defer mu.RUnlock()

	if statusQuery == "" {
		c.JSON(http.StatusOK, gin.H{"data": tasks})
		return
	}
	var filteredTasks []models.Task
	for _, task := range tasks {
		if string(task.Status) == statusQuery {
			filteredTasks = append(filteredTasks, task)
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": filteredTasks})

}

func PatchTask(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)

	var input models.PatchTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	defer mu.Unlock()

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == uint64(id) {
			if input.Title != nil {
				tasks[i].Title = *input.Title
			}
			if input.Description != nil {
				tasks[i].Description = *input.Description
			}
			if input.Status != nil {
				tasks[i].Status = *input.Status
			}
			c.JSON(http.StatusOK, gin.H{"data": tasks[i]})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
}
