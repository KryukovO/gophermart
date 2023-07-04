package handlers

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KryukovO/gophermart/internal/entities"
	mocks "github.com/KryukovO/gophermart/internal/mocks/repository"
	"github.com/KryukovO/gophermart/internal/usecases"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ErrRecorderIsNil = errors.New("empty response recorder")
	ErrURLIsEmpty    = errors.New("empty URL")
)

func newEchoContext(
	rec *httptest.ResponseRecorder,
	method, url string, body io.Reader,
) (echo.Context, error) {
	if rec == nil {
		return nil, ErrRecorderIsNil
	}

	if url == "" {
		return nil, ErrURLIsEmpty
	}

	e := echo.New()
	req := httptest.NewRequest(method, url, body)

	ctx := e.NewContext(req, rec)
	ctx.SetPath(url)

	return ctx, nil
}

func TestUserRequestRegisterHandler(t *testing.T) {
	var (
		secret = []byte("secret")
		url    = "/api/user/register"
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
				mock.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
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
				mock.EXPECT().Register(gomock.Any(), gomock.Any()).Return(entities.ErrUserAlreadyExists)
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := mocks.NewMockUserRepo(ctrl)

		if test.prepare != nil {
			test.prepare(mock)
		}

		rec := httptest.NewRecorder()
		echoCtx, err := newEchoContext(rec, http.MethodPost, url, bytes.NewReader(test.args.body))
		require.NoError(t, err)

		c := UserController{
			user:          usecases.NewUserUseCase(mock, time.Minute),
			secret:        secret,
			tokenLifetime: time.Minute,
		}
		err = c.userRequestHandler(c.user.Register)(echoCtx)
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
		url    = "/api/user/login"
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
				mock.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil)
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
				mock.EXPECT().Login(gomock.Any(), gomock.Any()).Return(entities.ErrInvalidLoginPassword)
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
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := mocks.NewMockUserRepo(ctrl)

		if test.prepare != nil {
			test.prepare(mock)
		}

		rec := httptest.NewRecorder()
		echoCtx, err := newEchoContext(rec, http.MethodPost, url, bytes.NewReader(test.args.body))
		require.NoError(t, err)

		c := UserController{
			user:          usecases.NewUserUseCase(mock, time.Minute),
			secret:        secret,
			tokenLifetime: time.Minute,
		}
		err = c.userRequestHandler(c.user.Login)(echoCtx)
		require.NoError(t, err)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, test.wants.status, res.StatusCode)

		if test.wants.setCookie {
			assert.NotEmpty(t, res.Cookies())
		}
	}
}
