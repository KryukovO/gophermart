package usecases

import (
	"context"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
)

type User interface {
	Register(ctx context.Context, user *entities.User, secret []byte) error
	Login(ctx context.Context, user *entities.User, secret []byte) error
}

type Order interface {
	Orders() error
	AddOrder() error
}

type Balance interface {
	Balance() error
	Withdraw() error
	Withdrawals() error
}
