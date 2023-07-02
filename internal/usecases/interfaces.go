package usecases

import (
	"context"

	"github.com/KryukovO/gophermart/internal/entities"
)

type Repo interface {
	Ping(ctx context.Context) error
	Close() error
}

type User interface {
	Register(ctx context.Context, user *entities.User, secret []byte) error
	Login(ctx context.Context, user *entities.User, secret []byte) error
}

type UserRepo interface {
	Register(ctx context.Context, user *entities.User) error
	Login(ctx context.Context, user *entities.User) error
}

type Order interface {
	Orders() error
	AddOrder() error
}

type OrderRepo interface {
	Orders() error
	AddOrder() error
}

type Balance interface {
	Balance() error
	Withdraw() error
	Withdrawals() error
}

type BalanceRepo interface {
	Balance() error
	Withdraw() error
	Withdrawals() error
}
