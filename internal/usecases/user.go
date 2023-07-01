package usecases

type UserUseCase struct {
	repo UserRepo
}

func NewUserUseCase(repo UserRepo) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (uc *UserUseCase) Register() error {
	return nil
}

func (uc *UserUseCase) Login() error {
	return nil
}

func (uc *UserUseCase) Orders() error {
	return nil
}

func (uc *UserUseCase) AddOrder() error {
	return nil
}

func (uc *UserUseCase) Balance() error {
	return nil
}

func (uc *UserUseCase) Withdraw() error {
	return nil
}

func (uc *UserUseCase) Withdrawals() error {
	return nil
}
