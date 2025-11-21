package data

import (
	"strconv"
	"task_manager/models"
)

var tasks = make(map[int]models.Task)
var nextID = 1

func init() {
	t1 := models.Task{
		ID:          nextID,
		Title:       "Finish homework",
		Description: "Math assignment",
		DueDate:     "2025-11-22",
		Status:      "pending",
	}
	tasks[nextID] = t1
	nextID++

	t2 := models.Task{
		ID:          nextID,
		Title:       "Buy groceries",
		Description: "Eggs, milk, bread",
		DueDate:     "2025-11-23",
		Status:      "ongoing",
	}
	tasks[nextID] = t2
	nextID++
}

// Get all tasks
func GetAllTasks() []models.Task {
	all := []models.Task{}
	for _, task := range tasks {
		all = append(all, task)
	}
	return all
}

// Get task by ID
func GetTaskByID(id string) models.Task {
	idInt, _ := strconv.Atoi(id)
	task, exists := tasks[idInt]
	if !exists {
		return models.Task{}
	}
	return task
}

// Create a new task
func CreateTask(task models.Task) models.Task {
	task.ID = nextID
	tasks[nextID] = task
	nextID++
	return task
}

// Update an existing task
func UpdateTask(id string, updatedTask models.Task) models.Task {
	idInt, _ := strconv.Atoi(id)
	_, exists := tasks[idInt]
	if !exists {
		return models.Task{}
	}
	updatedTask.ID = idInt
	tasks[idInt] = updatedTask
	return updatedTask
}

// Delete a task
func DeleteTask(id string) {
	idInt, _ := strconv.Atoi(id)
	delete(tasks, idInt)
}
