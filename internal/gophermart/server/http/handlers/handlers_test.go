package handlers

import (
	"testing"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/repository/mocks"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetHandlers(t *testing.T) {
	type args struct {
		server        *echo.Echo
		secret        []byte
		tokenLifetime time.Duration
		user          usecases.User
		order         usecases.Order
		balance       usecases.Balance
		logger        *log.Logger
	}

	type wants struct {
		wantErr bool
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "Correct setting",
			args: args{
				server:        echo.New(),
				secret:        []byte{},
				tokenLifetime: time.Second,
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				order:         usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				balance:       usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				logger:        log.New(),
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil logger",
			args: args{
				server:        echo.New(),
				secret:        []byte{},
				tokenLifetime: time.Second,
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				order:         usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				balance:       usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil server",
			args: args{
				secret:        []byte{},
				tokenLifetime: time.Second,
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				order:         usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				balance:       usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				logger:        log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
		{
			name: "Nil user",
			args: args{
				server:        echo.New(),
				secret:        []byte{},
				tokenLifetime: time.Second,
				order:         usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				balance:       usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				logger:        log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
		{
			name: "Nil order",
			args: args{
				server:        echo.New(),
				secret:        []byte{},
				tokenLifetime: time.Second,
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				balance:       usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				logger:        log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
		{
			name: "Nil balance",
			args: args{
				server:        echo.New(),
				secret:        []byte{},
				tokenLifetime: time.Second,
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				order:         usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				logger:        log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		err := SetHandlers(
			test.args.server, test.args.secret, test.args.tokenLifetime,
			test.args.user, test.args.order, test.args.balance, test.args.logger,
		)

		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
