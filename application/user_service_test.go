package application

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

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
	service := NewUserService(repo)

	err := service.Register("testuser", "password")
	assert.NoError(t, err)

	// Try to register the same user again
	err = service.Register("testuser", "password")
	assert.ErrorIs(t, err, ErrUsernameTaken)
}

func TestRegisterUser_EmptyUsername(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo)

	err := service.Register("", "password")
	assert.Error(t, err)
}

func TestRegisterUser_EmptyPassword(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo)

	err := service.Register("testuser", "")
	assert.Error(t, err)
}

func TestRegisterUser_PasswordIsHashed(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo)

	username := "secureuser"
	password := "supersecure"
	err := service.Register(username, password)
	assert.NoError(t, err)

	storedUser, err := repo.FindByUsername(username)
	assert.NoError(t, err)
	assert.NotEqual(t, password, storedUser.PasswordHash, "Password should be hashed")
}

func TestRegisterMultipleUsers(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo)

	users := []struct {
		username string
		password string
	}{
		{"alice", "password1"},
		{"bob", "password2"},
		{"charlie", "password3"},
	}

	for _, u := range users {
		err := service.Register(u.username, u.password)
		assert.NoError(t, err, "should register user %s", u.username)
	}

	for _, u := range users {
		_, err := repo.FindByUsername(u.username)
		assert.NoError(t, err, "should find user %s", u.username)
	}
}

func TestMockUserRepo_FindByUsername_NotFound(t *testing.T) {
	repo := newMockUserRepo()

	_, err := repo.FindByUsername("ghost")
	assert.Error(t, err, "should error for non-existent user")
}

func TestMockUserRepo_Create_Duplicate(t *testing.T) {
	repo := newMockUserRepo()
	user := &domain.User{
		ID:           uuid.New(),
		Username:     "duplicate",
		PasswordHash: "hashed",
	}

	err := repo.Create(user)
	assert.NoError(t, err, "first create should succeed")

	err = repo.Create(user)
	assert.Error(t, err, "duplicate create should fail")
}
