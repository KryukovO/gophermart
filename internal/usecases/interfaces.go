package usecases

type User interface {
	Register() error
	Login() error
	Orders() error
	AddOrder() error
	Balance() error
	Withdraw() error
	Withdrawals() error
}

type UserRepo interface {
	Close() error
	Register() error
	Login() error
	Orders() error
	AddOrder() error
	Balance() error
	Withdraw() error
	Withdrawals() error
}
