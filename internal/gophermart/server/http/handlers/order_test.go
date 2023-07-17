package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KryukovO/gophermart/internal/gophermart/entities"
	"github.com/KryukovO/gophermart/internal/gophermart/repository/mocks"
	"github.com/KryukovO/gophermart/internal/gophermart/usecases"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			logger: logrus.StandardLogger(),
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
			logger: logrus.StandardLogger(),
		}

		err := oc.ordersHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)
		assert.Equal(t, test.wants.contentType, res.Header.Get("Content-Type"))
	}
}
