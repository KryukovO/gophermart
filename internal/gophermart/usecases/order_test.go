package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	order1 := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "New",
		Accrual:    0,
		UploadedAt: time.Now(),
	}
	order2 := entities.Order{
		UserID:     2,
		Number:     "4561261212345464",
		Status:     "New",
		Accrual:    0,
		UploadedAt: time.Now(),
	}

	type args struct {
		order *entities.Order
	}

	type wants struct {
		wantErr bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockOrderRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct addition of the order",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().AddOrder(gomock.Any(), &order1).Return(nil)
			},
			args: args{
				order: &order1,
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Order has already been added",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().AddOrder(gomock.Any(), &order1).Return(entities.ErrOrderAlreadyAdded)
			},
			args: args{
				order: &order1,
			},
			wants: wants{
				wantErr: true,
			},
		},
		{
			name: "Order has already been added by other",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().AddOrder(gomock.Any(), &order1).Return(entities.ErrOrderAddedByOther)
			},
			args: args{
				order: &order1,
			},
			wants: wants{
				wantErr: true,
			},
		},
		{
			name: "Invalid order number",
			args: args{
				order: &order2,
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockOrderRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewOrderUseCase(repo, time.Minute)

		err := order.AddOrder(context.Background(), test.args.order)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestOrders(t *testing.T) {
	order1 := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "New",
		Accrual:    0,
		UploadedAt: time.Now(),
	}

	type args struct {
		userID int64
	}

	type wants struct {
		expected []entities.Order
		wantErr  bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockOrderRepo)
		args    args
		wants   wants
	}{
		{
			name: "Order list",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().Orders(gomock.Any(), gomock.Any()).Return([]entities.Order{order1}, nil)
			},
			args: args{
				userID: int64(1),
			},
			wants: wants{
				expected: []entities.Order{order1},
				wantErr:  false,
			},
		},
		{
			name: "Orders not found",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().Orders(gomock.Any(), gomock.Any()).Return([]entities.Order{}, nil)
			},
			args: args{
				userID: int64(1),
			},
			wants: wants{
				expected: []entities.Order{},
				wantErr:  false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockOrderRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewOrderUseCase(repo, time.Minute)

		orders, err := order.Orders(context.Background(), test.args.userID)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.wants.expected, orders)
		}
	}
}

func TestProcessableOrders(t *testing.T) {
	order1 := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "New",
		Accrual:    0,
		UploadedAt: time.Now(),
	}

	type wants struct {
		expected []entities.Order
		wantErr  bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockOrderRepo)
		wants   wants
	}{
		{
			name: "Order list",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().ProcessableOrders(gomock.Any()).Return([]entities.Order{order1}, nil)
			},
			wants: wants{
				expected: []entities.Order{order1},
				wantErr:  false,
			},
		},
		{
			name: "Orders not found",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().ProcessableOrders(gomock.Any()).Return([]entities.Order{}, nil)
			},
			wants: wants{
				expected: []entities.Order{},
				wantErr:  false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockOrderRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewOrderUseCase(repo, time.Minute)

		orders, err := order.ProcessableOrders(context.Background())
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.wants.expected, orders)
		}
	}
}

func TestUpdateOrder(t *testing.T) {
	order1 := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "PROCESSED",
		Accrual:    500,
		UploadedAt: time.Now().AddDate(0, 0, -1),
	}

	type wants struct {
		wantErr bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockOrderRepo)
		wants   wants
	}{
		{
			name: "Successful update",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().UpdateOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			wants: wants{
				wantErr: false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockOrderRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewOrderUseCase(repo, time.Minute)

		err := order.UpdateOrder(context.Background(), &order1)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
