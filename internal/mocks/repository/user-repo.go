// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KryukovO/gophermart/internal/usecases (interfaces: UserRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	entities "github.com/KryukovO/gophermart/internal/entities"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockUserRepo) Login(arg0 context.Context, arg1 *entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockUserRepoMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserRepo)(nil).Login), arg0, arg1)
}

// Register mocks base method.
func (m *MockUserRepo) Register(arg0 context.Context, arg1 *entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserRepoMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserRepo)(nil).Register), arg0, arg1)
}
