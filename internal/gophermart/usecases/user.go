package usecases

import (
	"context"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository"
)

type UserUseCase struct {
	repo    repository.UserRepo
	timeout time.Duration
}

func NewUserUseCase(repo repository.UserRepo, timeout time.Duration) *UserUseCase {
	return &UserUseCase{
		repo:    repo,
		timeout: timeout,
	}
}

func (uc *UserUseCase) Register(ctx context.Context, user *entities.User, secret []byte) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	err := user.Encrypt(secret)
	if err != nil {
		return err
	}

	return uc.repo.CreateUser(ctx, user)
}

func (uc *UserUseCase) Login(ctx context.Context, user *entities.User, secret []byte) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	err := uc.repo.UserByLogin(ctx, user)
	if err != nil {
		return err
	}

	return user.Validate(secret)
}
