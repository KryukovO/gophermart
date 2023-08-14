package accrualconnector

import (
	"context"
	"testing"
	"time"

	accmock "github.com/KryukovO/gophermart/internal/accrual/mocks"
	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository/mocks"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccrualConnector(t *testing.T) {
	type args struct {
		accrualAddr string
		workers     uint
		interval    time.Duration
		order       usecases.Order
		balance     usecases.Balance
		logger      *log.Logger
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Correct creation",
			args: args{
				accrualAddr: "http://localhost:8080",
				workers:     3,
				interval:    time.Second,
				order:       usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				balance:     usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				logger:      log.New(),
			},
		},
		{
			name: "Nil logger",
			args: args{
				accrualAddr: "http://localhost:8080",
				workers:     3,
				interval:    time.Second,
				order:       mocks.NewMockOrderRepo(gomock.NewController(t)),
				balance:     mocks.NewMockBalanceRepo(gomock.NewController(t)),
				logger:      log.New(),
			},
		},
	}

	for _, test := range tests {
		con := NewAccrualConnector(
			test.args.accrualAddr, test.args.workers, test.args.interval,
			test.args.order, test.args.balance, test.args.logger,
		)

		require.NotNil(t, con)

		assert.Equal(t, test.args.accrualAddr, con.accrualAddr)
		assert.Equal(t, test.args.workers, con.workers)
		assert.Equal(t, test.args.interval, con.interval)
		assert.Equal(t, test.args.order, con.order)
		assert.Equal(t, test.args.balance, con.balance)

		if test.args.logger != nil {
			assert.Equal(t, test.args.logger, con.logger)
		} else {
			assert.NotNil(t, con.logger)
		}
	}
}

func TestRun(t *testing.T) {
	ts := time.Now().AddDate(0, 0, 1)
	order := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "NEW",
		UploadedAt: ts,
	}
	orderExpected := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "PROCESSED",
		Accrual:    500,
		UploadedAt: ts,
	}
	balance := entities.BalanceChange{
		UserID:    1,
		Operation: "refill",
		Order:     "4561261212345467",
		Sum:       500,
	}

	type args struct {
		timeout  time.Duration
		shutdown bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Shutdown by context",
			args: args{
				timeout: 3 * time.Second,
			},
		},
		{
			name: "Shutdown by command",
			args: args{
				shutdown: true,
			},
		},
	}

	accrual := accmock.NewMockAccrual()
	defer accrual.Close()

	for _, test := range tests {
		ctr := gomock.NewController(t)
		orderRepo := mocks.NewMockOrderRepo(ctr)
		balanceRepo := mocks.NewMockBalanceRepo(ctr)

		orderRepo.EXPECT().ProcessableOrders(gomock.Any()).AnyTimes().Return([]entities.Order{order}, nil)
		orderRepo.EXPECT().UpdateOrder(gomock.Any(), mocks.OrderMatcher(&orderExpected)).AnyTimes().Return(nil)
		balanceRepo.EXPECT().ChangeBalance(gomock.Any(), mocks.BalanceChangeMatcher(&balance)).AnyTimes().Return(nil)

		con := NewAccrualConnector(
			accrual.URL, 1, time.Second,
			usecases.NewOrderUseCase(orderRepo, time.Second),
			usecases.NewBalanceUseCase(balanceRepo, time.Second),
			nil,
		)

		ctx, cancel := context.WithCancel(context.Background())
		if !test.args.shutdown {
			ctx, cancel = context.WithTimeout(context.Background(), test.args.timeout)
		}

		go con.Run(ctx)

		if test.args.shutdown {
			con.Shutdown(context.Background())
		} else {
			<-ctx.Done()
		}

		cancel()
	}
}

func TestGenerateOrderTasks(t *testing.T) {
	orders := []entities.Order{
		{
			UserID:     1,
			Number:     "4561261212345467",
			Status:     "NEW",
			UploadedAt: time.Now().AddDate(0, 0, -1),
		},
		{
			UserID:     2,
			Number:     "4861261212345464",
			Status:     "NEW",
			UploadedAt: time.Now().AddDate(0, 0, -1),
		},
	}

	con := AccrualConnector{workers: 3}
	ch := con.generateOrderTasks(context.Background(), orders)
	tasks := make([]entities.Order, 0, len(orders))

	for task := range ch {
		tasks = append(tasks, task)
	}

	assert.ElementsMatch(t, orders, tasks)
}

func TestOrderTaskWorker(t *testing.T) {
	ts := time.Now().AddDate(0, 0, 1)
	order := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "NEW",
		UploadedAt: ts,
	}
	orderExpected := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "PROCESSED",
		Accrual:    500,
		UploadedAt: ts,
	}
	balance := entities.BalanceChange{
		UserID:    1,
		Operation: "refill",
		Order:     "4561261212345467",
		Sum:       500,
	}

	accrual := accmock.NewMockAccrual()
	defer accrual.Close()

	ctr := gomock.NewController(t)
	orderRepo := mocks.NewMockOrderRepo(ctr)
	balanceRepo := mocks.NewMockBalanceRepo(ctr)

	orderRepo.EXPECT().UpdateOrder(gomock.Any(), mocks.OrderMatcher(&orderExpected)).Return(nil)
	balanceRepo.EXPECT().ChangeBalance(gomock.Any(), mocks.BalanceChangeMatcher(&balance)).Return(nil)

	con := AccrualConnector{
		accrualAddr: accrual.URL,
		order:       orderRepo,
		balance:     balanceRepo,
	}

	ch := make(chan entities.Order, 1)
	ch <- order
	close(ch)

	err := con.orderTaskWorker(context.Background(), ch)

	assert.NoError(t, err)
}

func TestDoRequest(t *testing.T) {
	order := entities.AccrualOrder{
		Order:   "4561261212345467",
		Status:  "PROCESSED",
		Accrual: 500,
	}

	accrual := accmock.NewMockAccrual()
	defer accrual.Close()

	con := AccrualConnector{accrualAddr: accrual.URL}
	res, err := con.doRequest(context.Background(), accrual.Client(), order.Order)

	assert.NoError(t, err)
	assert.Equal(t, order, res)
}
