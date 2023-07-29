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

func TestNewOrderController(t *testing.T) {
	type args struct {
		order     usecases.Order
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
				order:     usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
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
				order:     usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
				mwManager: middleware.NewManager([]byte("secret"), log.New()),
				logger:    nil,
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil order",
			args: args{
				order:     nil,
				mwManager: middleware.NewManager([]byte("secret"), log.New()),
				logger:    log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		ctrl, err := NewOrderController(
			test.args.order, test.args.mwManager, test.args.logger,
		)

		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, ctrl)

			assert.Equal(t, test.args.order, ctrl.order)
			assert.Equal(t, test.args.mwManager, ctrl.mw)

			if test.args.logger != nil {
				assert.Equal(t, test.args.logger, ctrl.logger)
			} else {
				assert.NotNil(t, ctrl.logger)
			}
		}
	}
}

func TestOrderMapHandlers(t *testing.T) {
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
		ctrl, err := NewOrderController(
			usecases.NewOrderUseCase(mocks.NewMockOrderRepo(gomock.NewController(t)), time.Second),
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

func TestAddOrderHandler(t *testing.T) {
	path := "/api/user/orders"

	type args struct {
		userID interface{}
		body   []byte
	}

	type wants struct {
		status int
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
				mock.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				userID: int64(1),
				body:   []byte("4561261212345467"),
			},
			wants: wants{
				status: http.StatusAccepted,
			},
		},
		{
			name: "Order has already been added",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(entities.ErrOrderAlreadyAdded)
			},
			args: args{
				userID: int64(1),
				body:   []byte("4561261212345467"),
			},
			wants: wants{
				status: http.StatusOK,
			},
		},
		{
			name: "Order has already been added by other",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(entities.ErrOrderAddedByOther)
			},
			args: args{
				userID: int64(1),
				body:   []byte("4561261212345467"),
			},
			wants: wants{
				status: http.StatusConflict,
			},
		},
		{
			name: "User unauthorized",
			args: args{
				body: []byte("4561261212345467"),
			},
			wants: wants{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "Invalid order number",
			args: args{
				userID: int64(1),
				body:   []byte("4561261212345464"),
			},
			wants: wants{
				status: http.StatusUnprocessableEntity,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockOrderRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(test.args.body))
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)
		echoCtx.Set("userID", test.args.userID)

		oc := OrderController{
			order:  usecases.NewOrderUseCase(repo, time.Minute),
			logger: log.StandardLogger(),
		}

		err := oc.addOrderHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)
	}
}

func TestOrdersHandler(t *testing.T) {
	path := "/api/user/orders"
	order := entities.Order{
		UserID:     1,
		Number:     "4561261212345467",
		Status:     "New",
		Accrual:    0,
		UploadedAt: time.Now(),
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
		prepare func(mock *mocks.MockOrderRepo)
		args    args
		wants   wants
	}{
		{
			name: "Order list",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().Orders(gomock.Any(), gomock.Any()).Return([]entities.Order{order}, nil)
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
			name: "Orders not found",
			prepare: func(mock *mocks.MockOrderRepo) {
				mock.EXPECT().Orders(gomock.Any(), gomock.Any()).Return([]entities.Order{}, nil)
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
		repo := mocks.NewMockOrderRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)
		echoCtx.Set("userID", test.args.userID)

		oc := OrderController{
			order:  usecases.NewOrderUseCase(repo, time.Minute),
			logger: log.StandardLogger(),
		}

		err := oc.ordersHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)
		assert.Equal(t, test.wants.contentType, res.Header.Get("Content-Type"))
	}
}
