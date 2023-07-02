package usecases

import "time"

type UserUseCase struct {
	repo    UserRepo
	timeout time.Duration
}

func NewUserUseCase(repo UserRepo, timeout time.Duration) *UserUseCase {
	return &UserUseCase{
		repo:    repo,
		timeout: timeout,
	}
}

func (uc *UserUseCase) CreateUser() error {
	return nil
}

func (uc *UserUseCase) Login() error {
	return nil
}
