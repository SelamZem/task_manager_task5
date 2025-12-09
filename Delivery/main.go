package main

import (
	"context"
	"log"
	"task_manager/Delivery/controllers"
	repositories "task_manager/Repositories"
	usecases "task_manager/Usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Initialize MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	// Wire repository using MongoTaskRepository
	var repo repositories.TaskRepository = repositories.NewMongoTaskRepository(client, "task_manager", "tasks")

	// Setup HTTP server
	r := gin.Default()
	taskGroup := r.Group("/tasks")
	uc := usecases.TaskUsecase{Repo: repo}
	controllers.RegisterTaskRoutes(taskGroup, uc)
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
