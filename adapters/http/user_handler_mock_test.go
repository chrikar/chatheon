package http

import (
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) Register(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *mockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}
