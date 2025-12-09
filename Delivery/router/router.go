package router

import (
	"task_manager/Delivery/controllers"
	repositories "task_manager/Repositories"
	usecases "task_manager/Usecases"

	"github.com/gin-gonic/gin"
)

func SetupRouter(repo repositories.TaskRepository) *gin.Engine {
	r := gin.Default()

	api := r.Group("/tasks")
	uc := usecases.TaskUsecase{Repo: repo}
	controllers.RegisterTaskRoutes(api, uc)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
