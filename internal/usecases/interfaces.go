package usecases

import "context"

type Repo interface {
	Ping(ctx context.Context) error
	Close() error
}

type User interface {
	CreateUser() error
	Login() error
}

type UserRepo interface {
	CreateUser() error
	Login() error
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
