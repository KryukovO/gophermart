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

func TestRegister(t *testing.T) {
	var (
		user1 = entities.User{
			Login:    "user1",
			Password: "1234",
		}
		secret = []byte("secret")
	)

	type args struct {
		user *entities.User
	}

	tests := []struct {
		name    string
		prepare func(mock *mocks.MockUserRepo)
		args    args
		wantErr bool
	}{
		{
			name: "Correct registration",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().AddUser(gomock.Any(), &user1).Return(nil)
			},
			args: args{
				user: &user1,
			},
			wantErr: false,
		},
		{
			name: "User already exists",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().AddUser(gomock.Any(), &user1).Return(entities.ErrUserAlreadyExists)
			},
			args: args{
				user: &user1,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockUserRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		user := NewUserUseCase(repo, time.Minute)

		err := user.Register(context.Background(), test.args.user, secret)
		if test.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

func TestLogin(t *testing.T) {
	var (
		user1 = entities.User{
			Login:    "user1",
			Password: "1234",
		}
		secret = []byte("secret")
	)

	type args struct {
		user *entities.User
	}

	type wants struct {
		wantErr bool
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
				user: &user1,
			},
			wants: wants{
				wantErr: false,
			},
		},
		{
			name: "Invalid login/password",
			prepare: func(mock *mocks.MockUserRepo) {
				mock.EXPECT().User(gomock.Any(), gomock.Any()).Return(entities.ErrInvalidLoginPassword)
			},
			args: args{
				user: &user1,
			},
			wants: wants{
				wantErr: true,
			},
		},
	}

	for _, test := range tests {
		repo := mocks.NewMockUserRepo(gomock.NewController(t))

		if test.prepare != nil {
			test.prepare(repo)
		}

		user := NewUserUseCase(repo, time.Minute)

		err := user.Login(context.Background(), test.args.user, secret)
		if test.wants.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
