package http

import (
	"github.com/chrikar/chatheon/application/ports"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

var _ ports.UserService = (*mockUserService)(nil) // Optional compile-time check

func (m *mockUserService) Register(username, password string) error {
	args := m.Called(username, password)
	return args.Error(0)
}

func (m *mockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}
