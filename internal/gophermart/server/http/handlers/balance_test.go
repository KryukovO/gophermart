package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository/mocks"
	"github.com/KryukovO/gophermart/internal/gophermart/server/http/middleware"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBalanceController(t *testing.T) {
	type args struct {
		balance   usecases.Balance
		mwManager *middleware.Manager
		logger    *log.Logger
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
			name: "Correct creation",
			args: args{
				balance:   usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				mwManager: middleware.NewManager([]byte("secret"), log.New()),
				logger:    log.New(),
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil logger",
			args: args{
				balance:   usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
				mwManager: middleware.NewManager([]byte("secret"), log.New()),
				logger:    nil,
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil balance",
			args: args{
				balance:   nil,
				mwManager: middleware.NewManager([]byte("secret"), log.New()),
				logger:    log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		ctrl, err := NewBalanceController(
			test.args.balance, test.args.mwManager, test.args.logger,
		)

		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, ctrl)

			assert.Equal(t, test.args.balance, ctrl.balance)
			assert.Equal(t, test.args.mwManager, ctrl.mw)

			if test.args.logger != nil {
				assert.Equal(t, test.args.logger, ctrl.logger)
			} else {
				assert.NotNil(t, ctrl.logger)
			}
		}
	}
}

func TestBalanceMapHandlers(t *testing.T) {
	type args struct {
		group *echo.Group
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
			name: "Correct mapping",
			args: args{
				group: echo.New().Group("/"),
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil group",
			args: args{
				group: nil,
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		ctrl, err := NewBalanceController(
			usecases.NewBalanceUseCase(mocks.NewMockBalanceRepo(gomock.NewController(t)), time.Second),
			middleware.NewManager([]byte("secret"), log.New()),
			log.New(),
		)

		require.NoError(t, err)
		require.NotNil(t, ctrl)

		err = ctrl.MapHandlers(test.args.group)

		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestBalanceHandler(t *testing.T) {
	path := "/api/user/balance"
	balance := entities.Balance{
		UserID:    1,
		Current:   500,
		Withdrawn: 42,
	}

	type args struct {
		userID interface{}
	}

	type wants struct {
		status      int
		contentType string
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
				status:      http.StatusOK,
				contentType: "application/json; charset=UTF-8",
			},
		},
		{
			name: "User unauthorized",
			args: args{},
			wants: wants{
				status: http.StatusUnauthorized,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockBalanceRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)
		echoCtx.Set("userID", test.args.userID)

		bc := BalanceController{
			balance: usecases.NewBalanceUseCase(repo, time.Minute),
		}

		err := bc.balanceHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)
		assert.Equal(t, test.wants.contentType, res.Header.Get("Content-Type"))
	}
}

func TestWithdrawHandler(t *testing.T) {
	path := "/api/user/balance/withdraw"

	type args struct {
		userID interface{}
		body   []byte
	}

	type wants struct {
		status int
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockBalanceRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct withdraw request",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().ChangeBalance(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				userID: int64(1),
				body:   []byte(`{"order":"4561261212345467","sum":751}`),
			},
			wants: wants{
				status: http.StatusOK,
			},
		},
		{
			name: "Not enough funds",
			prepare: func(mock *mocks.MockBalanceRepo) {
				mock.EXPECT().ChangeBalance(gomock.Any(), gomock.Any()).Return(entities.ErrNotEnoughFunds)
			},
			args: args{
				body:   []byte(`{"order":"4561261212345467","sum":751}`),
				userID: int64(1),
			},
			wants: wants{
				status: http.StatusPaymentRequired,
			},
		},
		{
			name: "Invalid order number",
			args: args{
				body:   []byte(`{"order":"4561261212345464","sum":751}`),
				userID: int64(1),
			},
			wants: wants{
				status: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Incorrect request body #1",
			args: args{
				body:   []byte(`{"order":"4561261212345467","sum":true}`),
				userID: int64(1),
			},
			wants: wants{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Incorrect request body #2",
			args: args{
				body:   []byte(`{"order":"4561261212345467","sum":-751}`),
				userID: int64(1),
			},
			wants: wants{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Incorrect request body #3",
			args: args{
				body:   []byte(`{"order":"4561261212345467","sum":0}`),
				userID: int64(1),
			},
			wants: wants{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "User unauthorized",
			args: args{},
			wants: wants{
				status: http.StatusUnauthorized,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockBalanceRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(test.args.body))
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)
		echoCtx.Set("userID", test.args.userID)

		bc := BalanceController{
			balance: usecases.NewBalanceUseCase(repo, time.Minute),
			logger:  log.StandardLogger(),
		}

		err := bc.withdrawHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)
	}
}

func TestWithdrawalsHandler(t *testing.T) {
	path := "/api/user/withdrawals"
	change := entities.BalanceChange{
		UserID:      1,
		Order:       "2377225624",
		Sum:         751,
		ProcessedAt: time.Now(),
	}

	type args struct {
		userID interface{}
	}

	type wants struct {
		status      int
		contentType string
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
				status:      http.StatusOK,
				contentType: "application/json; charset=UTF-8",
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
				status: http.StatusNoContent,
			},
		},
		{
			name: "User unauthorized",
			args: args{},
			wants: wants{
				status: http.StatusUnauthorized,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockBalanceRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)
		echoCtx.Set("userID", test.args.userID)

		bc := BalanceController{
			balance: usecases.NewBalanceUseCase(repo, time.Minute),
			logger:  log.StandardLogger(),
		}

		err := bc.withdrawalsHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)
		assert.Equal(t, test.wants.contentType, res.Header.Get("Content-Type"))
	}
}
