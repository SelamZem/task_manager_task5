package usecases

import (
	domain "task_manager/Domain"
	repositories "task_manager/Repositories"
)

type UserUsecase struct {
	Repo repositories.UserRepository
}

func (u *UserUsecase) Register(user domain.User) (domain.User, error) {
	return u.Repo.Create(user)
}

func (u *UserUsecase) GetByID(id string) (domain.User, error) {
	return u.Repo.GetByID(id)
}

func (u *UserUsecase) GetByUsername(username string) (domain.User, error) {
	return u.Repo.GetByUsername(username)
}

func (u *UserUsecase) PromoteToAdmin(username string) (domain.User, error) {
	return u.Repo.PromoteToAdmin(username)
}
