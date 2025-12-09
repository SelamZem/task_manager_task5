package usecases

import (
	domain "task_manager/Domain"
	repositories "task_manager/Repositories"
)

type TaskUsecase struct {
	Repo repositories.TaskRepository
}

func (u *TaskUsecase) GetAllTasks() ([]domain.Task, error) {
	return u.Repo.GetAllTasks()
}

func (u *TaskUsecase) GetTaskByID(id string) (domain.Task, error) {
	return u.Repo.GetTaskByID(id)
}

func (u *TaskUsecase) CreateTask(task domain.Task) (domain.Task, error) {
	return u.Repo.CreateTask(task)
}

func (u *TaskUsecase) UpdateTask(id string, task domain.Task) (domain.Task, error) {
	return u.Repo.UpdateTask(id, task)
}

func (u *TaskUsecase) DeleteTask(id string) error {
	return u.Repo.DeleteTask(id)
}
