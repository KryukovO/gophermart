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
	"github.com/labstack/echo"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRequestRegisterHandler(t *testing.T) {
	var (
		secret = []byte("secret")
		path   = "/api/user/register"
	)

	type args struct {
		body []byte
	}

	type wants struct {
		status    int
		setCookie bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockUserRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct registration",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				body: []byte(`{"login":"user1","password":"1234"}`),
			},
			wants: wants{
				status:    http.StatusOK,
				setCookie: true,
			},
		},
		{
			name: "User already exists",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(entities.ErrUserAlreadyExists)
			},
			args: args{
				body: []byte(`{"login":"user1","password":"1234"}`),
			},
			wants: wants{
				status:    http.StatusConflict,
				setCookie: false,
			},
		},
		{
			name:    "Incorrect request body",
			prepare: nil,
			args: args{
				body: []byte(`{"login":"user1","password":true}`),
			},
			wants: wants{
				status:    http.StatusBadRequest,
				setCookie: false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockUserRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(test.args.body))
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)

		c := UserController{
			user:          usecases.NewUserUseCase(repo, time.Minute),
			secret:        secret,
			tokenLifetime: time.Minute,
		}
		err := c.userRequestHandler(c.user.Register)(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)

		if test.wants.setCookie {
			assert.NotEmpty(t, res.Cookies())
		}
	}
}

func TestUserRequestLoginHandler(t *testing.T) {
	var (
		secret = []byte("secret")
		path   = "/api/user/login"
	)

	type args struct {
		body []byte
	}

	type wants struct {
		status    int
		setCookie bool
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockUserRepo)
		args    args
		wants   wants
	}{
		{
			name: "Correct login",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().User(gomock.Any(), gomock.Any()).Return(nil)
			},
			args: args{
				body: []byte(`{"login":"user1","password":"1234"}`),
			},
			wants: wants{
				status:    http.StatusOK,
				setCookie: true,
			},
		},
		{
			name: "Invalid login/password",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().User(gomock.Any(), gomock.Any()).Return(entities.ErrInvalidLoginPassword)
			},
			args: args{
				body: []byte(`{"login":"user1","password":"1234"}`),
			},
			wants: wants{
				status:    http.StatusUnauthorized,
				setCookie: false,
			},
		},
		{
			name:    "Incorrect request body",
			prepare: nil,
			args: args{
				body: []byte(`{"login":"user1","password":true}`),
			},
			wants: wants{
				status:    http.StatusBadRequest,
				setCookie: false,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockUserRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, path, bytes.NewBuffer(test.args.body))
		server := echo.New()
		echoCtx := server.NewContext(req, rec)

		echoCtx.SetPath(path)

		c := UserController{
			user:          usecases.NewUserUseCase(repo, time.Minute),
			secret:        secret,
			tokenLifetime: time.Minute,
		}
		err := c.userRequestHandler(c.user.Login)(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)

		if test.wants.setCookie {
			assert.NotEmpty(t, res.Cookies())
		}
	}
}
