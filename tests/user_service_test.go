package tests

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/application"
	"github.com/chrikar/chatheon/domain"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(user *domain.User) error {
	if _, exists := m.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	m.users[user.Username] = user
	return nil
}

func (m *mockUserRepo) FindByUsername(username string) (*domain.User, error) {
	user, exists := m.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func TestRegisterUser(t *testing.T) {
	repo := newMockUserRepo()
	service := application.NewUserService(repo)

	err := service.Register("testuser", "password")
	assert.NoError(t, err)

	// Try to register the same user again
	err = service.Register("testuser", "password")
	assert.ErrorIs(t, err, application.ErrUsernameTaken)
}