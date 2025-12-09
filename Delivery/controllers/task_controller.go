package controllers

import (
	"net/http"
	domain "task_manager/Domain"
	usecases "task_manager/Usecases"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(rg *gin.RouterGroup, uc usecases.TaskUsecase) {
	rg.GET("/", func(c *gin.Context) { listTasks(c, uc) })
	rg.GET("/:id", func(c *gin.Context) { getTask(c, uc) })
	rg.POST("/", func(c *gin.Context) { createTask(c, uc) })
	rg.PUT("/:id", func(c *gin.Context) { updateTask(c, uc) })
	rg.DELETE("/:id", func(c *gin.Context) { deleteTask(c, uc) })
}

func listTasks(c *gin.Context, uc usecases.TaskUsecase) {
	tasks, err := uc.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context, uc usecases.TaskUsecase) {
	id := c.Param("id")
	task, err := uc.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func createTask(c *gin.Context, uc usecases.TaskUsecase) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	created, err := uc.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func updateTask(c *gin.Context, uc usecases.TaskUsecase) {
	id := c.Param("id")
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	updated, err := uc.UpdateTask(id, task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func deleteTask(c *gin.Context, uc usecases.TaskUsecase) {
	id := c.Param("id")
	err := uc.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
