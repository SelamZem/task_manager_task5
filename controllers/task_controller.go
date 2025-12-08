package controllers

import (
	"net/http"
	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(rg *gin.RouterGroup) {
	rg.GET("/", listTasks)
	rg.GET("/:id", getTask)
	rg.POST("/", createTask)
	rg.PUT("/:id", updateTask)
	rg.DELETE("/:id", deleteTask)
}

func listTasks(c *gin.Context) {
	tasks := data.GetAllTasks()
	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context) {
	id := c.Param("id")
	task := data.GetTaskByID(id)

	if task.ID.IsZero() { // zero-value indicates not found
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)

}

func createTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"an error occured loading ur data": "invalid JSON"})
		return
	}

	created := data.CreateTask(task)
	c.JSON(http.StatusCreated, created)

}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"an error occured loading ur data": "invalid JSON"})
		return
	}

	updated := data.UpdateTask(id, task)
	if updated.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"an error occured": "task not found"})
		return
	}

	c.JSON(http.StatusOK, updated)

}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	data.DeleteTask(id)
	c.Status(http.StatusNoContent)
}
