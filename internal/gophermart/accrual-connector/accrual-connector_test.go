package accrualconnector

import (
	"context"
	"testing"
	"time"

	accmock "github.com/KryukovO/gophermart/internal/accrual/mocks"
	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
