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
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserController(t *testing.T) {
	type args struct {
		user          usecases.User
		secret        []byte
		tokenLifetime time.Duration
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
			name: "Correct creation",
			args: args{
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				secret:        []byte{},
				tokenLifetime: time.Second,
				logger:        log.New(),
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil logger",
			args: args{
				user:          usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
				secret:        []byte{},
				tokenLifetime: time.Second,
				logger:        nil,
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Nil user",
			args: args{
				user:          nil,
				secret:        []byte{},
				tokenLifetime: time.Second,
				logger:        log.New(),
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		ctrl, err := NewUserController(
			test.args.user, test.args.secret, test.args.tokenLifetime, test.args.logger,
		)

		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, ctrl)

			assert.Equal(t, test.args.user, ctrl.user)
			assert.Equal(t, test.args.secret, ctrl.secret)
			assert.Equal(t, test.args.tokenLifetime, ctrl.tokenLifetime)

			if test.args.logger != nil {
				assert.Equal(t, test.args.logger, ctrl.logger)
			} else {
				assert.NotNil(t, ctrl.logger)
			}
		}
	}
}

func TestUserMapHandlers(t *testing.T) {
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
		ctrl, err := NewUserController(
			usecases.NewUserUseCase(mocks.NewMockUserRepo(gomock.NewController(t)), time.Second),
			[]byte{},
			time.Second,
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

func TestRegisterHandler(t *testing.T) {
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

		uc := UserController{
			user:          usecases.NewUserUseCase(repo, time.Minute),
			secret:        secret,
			tokenLifetime: time.Minute,
			logger:        log.StandardLogger(),
		}
		err := uc.registerHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)

		if test.wants.setCookie {
			assert.NotEmpty(t, res.Cookies())
		}
	}
}

func TestLoginHandler(t *testing.T) {
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

		uc := UserController{
			user:          usecases.NewUserUseCase(repo, time.Minute),
			secret:        secret,
			tokenLifetime: time.Minute,
			logger:        log.StandardLogger(),
		}
		err := uc.loginHandler(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)

		if test.wants.setCookie {
			assert.NotEmpty(t, res.Cookies())
		}
	}
}
