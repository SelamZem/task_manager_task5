package main

import (
	"task_manager/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	taskGroup := r.Group("/tasks")
	controllers.RegisterTaskRoutes(taskGroup)

	r.Run(":8080")
}
