package data

import (
	"context"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var taskCollection *mongo.Collection
var ctx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI(getMongoURI())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		panic(err)
	}
	taskCollection = client.Database("task_manager").Collection("tasks")
}

// Get all tasks
func GetAllTasks() []models.Task {
	cursor, err := taskCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	for cursor.Next(ctx) {
		var t models.Task
		if err := cursor.Decode(&t); err == nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

// Get task by ID
func GetTaskByID(id string) models.Task {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}
	}
	var task models.Task
	err = taskCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		return models.Task{}
	}
	return task
}

// Create a new task
func CreateTask(task models.Task) models.Task {
	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}
	_, err := taskCollection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}
	}
	return task
}

// Update an existing task
func UpdateTask(id string, updatedTask models.Task) models.Task {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}
	}
	updatedTask.ID = objID
	_, err = taskCollection.ReplaceOne(ctx, bson.M{"_id": objID}, updatedTask)
	if err != nil {
		return models.Task{}
	}
	return updatedTask
}

// Delete a task
func DeleteTask(id string) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	taskCollection.DeleteOne(ctx, bson.M{"_id": objID})
}

// memory

// package data

// import (
// 	"strconv"
// 	"task_manager/models"
// )

// var tasks = make(map[int]models.Task)
// var nextID = 1

// func init() {
// 	t1 := models.Task{
// 		ID:          nextID,
// 		Title:       "Finish homework",
// 		Description: "Math assignment",
// 		DueDate:     "2025-11-22",
// 		Status:      "pending",
// 	}
// 	tasks[nextID] = t1
// 	nextID++

// 	t2 := models.Task{
// 		ID:          nextID,
// 		Title:       "Buy groceries",
// 		Description: "Eggs, milk, bread",
// 		DueDate:     "2025-11-23",
// 		Status:      "ongoing",
// 	}
// 	tasks[nextID] = t2
// 	nextID++
// }

// // Get all tasks
// func GetAllTasks() []models.Task {
// 	all := []models.Task{}
// 	for _, task := range tasks {
// 		all = append(all, task)
// 	}
// 	return all
// }

// // Get task by ID
// func GetTaskByID(id string) models.Task {
// 	idInt, _ := strconv.Atoi(id)
// 	task, exists := tasks[idInt]
// 	if !exists {
// 		return models.Task{}
// 	}
// 	return task
// }

// // Create a new task
// func CreateTask(task models.Task) models.Task {
// 	task.ID = nextID
// 	tasks[nextID] = task
// 	nextID++
// 	return task
// }

// // Update an existing task
// func UpdateTask(id string, updatedTask models.Task) models.Task {
// 	idInt, _ := strconv.Atoi(id)
// 	_, exists := tasks[idInt]
// 	if !exists {
// 		return models.Task{}
// 	}
// 	updatedTask.ID = idInt
// 	tasks[idInt] = updatedTask
// 	return updatedTask
// }

// // Delete a task
// func DeleteTask(id string) {
// 	idInt, _ := strconv.Atoi(id)
// 	delete(tasks, idInt)
// }
