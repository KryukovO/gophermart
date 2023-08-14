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

func TestBalance(t *testing.T) {
	balance := entities.Balance{
		UserID:    1,
		Current:   1000,
		Withdrawn: 500,
	}

	type args struct {
		userID int64
	}

	type wants struct {
		expected entities.Balance
		wantErr  bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockBalanceRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct balance request",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().Balance(gomock.Any(), gomock.Any()).Return(balance, nil)
			},
			args: args{
				userID: int64(1),
			},
			wants: wants{
				expected: balance,
				wantErr:  false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockBalanceRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewBalanceUseCase(repo, time.Minute)

		result, err := order.Balance(context.Background(), test.args.userID)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.wants.expected, result)
		}
	}
}

func TestChangeBalance(t *testing.T) {
	change1 := entities.BalanceChange{
		UserID:    1,
		Operation: entities.BalanceOperationWithdrawal,
		Order:     "4561261212345467",
		Sum:       1000,
	}
	change2 := entities.BalanceChange{
		UserID:    1,
		Operation: entities.BalanceOperationWithdrawal,
		Order:     "4561261212345464",
		Sum:       1000,
	}

	type args struct {
		change entities.BalanceChange
	}

	type wants struct {
		wantErr bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockBalanceRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct balance change",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().ChangeBalance(gomock.Any(), &change1).Return(nil)
			},
			args: args{
				change: change1,
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Invalid order number",
			args: args{
				change: change2,
			},
			wants: wants{
				wantErr: true,
			},
		},
		{
			name: "Not enough funds",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().ChangeBalance(gomock.Any(), &change1).Return(entities.ErrNotEnoughFunds)
			},
			args: args{
				change: change1,
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockBalanceRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewBalanceUseCase(repo, time.Minute)

		err := order.ChangeBalance(context.Background(), &test.args.change)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestWithdrawals(t *testing.T) {
	change := entities.BalanceChange{
		UserID:    1,
		Operation: entities.BalanceOperationWithdrawal,
		Order:     "4561261212345467",
		Sum:       1000,
	}

	type args struct {
		userID int64
	}

	type wants struct {
		expected []entities.BalanceChange
		wantErr  bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockBalanceRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct withdrawals request",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().Withdrawals(gomock.Any(), gomock.Any()).Return([]entities.BalanceChange{change}, nil)
			},
			args: args{
				userID: int64(1),
			},
			wants: wants{
				expected: []entities.BalanceChange{change},
				wantErr:  false,
			},
		},
		{
			name: "Withdrawals not found",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().Withdrawals(gomock.Any(), gomock.Any()).Return([]entities.BalanceChange{}, nil)
			},
			args: args{
				userID: int64(1),
			},
			wants: wants{
				expected: []entities.BalanceChange{},
				wantErr:  false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockBalanceRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		order := NewBalanceUseCase(repo, time.Minute)

		result, err := order.Withdrawals(context.Background(), test.args.userID)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.ElementsMatch(t, test.wants.expected, result)
		}
	}
}
