package usecases

import (
	"context"
	"math/rand"
	"time"

	"github.com/KryukovO/gophermart/internal/entities"
	"github.com/KryukovO/gophermart/internal/utils"
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

	salt, err := utils.GenerateRandomSalt(rand.NewSource(time.Now().UnixNano()))
	if err != nil {
		return err
	}

	user.Salt = salt

	err = user.Encrypt(secret)
	if err != nil {
		return err
	}

	return uc.repo.Register(ctx, user)
}

func (uc *UserUseCase) Login(ctx context.Context, user *entities.User, secret []byte) error {
	ctx, cancel := context.WithTimeout(ctx, uc.timeout)
	defer cancel()

	password := user.Password

	err := uc.repo.Login(ctx, user)
	if err != nil {
		return err
	}

	return user.Validate(password, secret)
}
