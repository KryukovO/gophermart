package usecases

import (
	"context"
	"time"

	"github.com/KryukovO/gophermart/internal/entities"
)

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

func (uc *UserUseCase) Register(ctx context.Context, user *entities.User, secret []byte) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	err := user.Encrypt(secret)
	if err != nil {
		return err
	}

	return uc.repo.Register(ctx, user)
}

func (uc *UserUseCase) Login(ctx context.Context, user *entities.User, secret []byte) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	err := uc.repo.Login(ctx, user)
	if err != nil {
		return err
	}

	return user.Validate(secret)
}
