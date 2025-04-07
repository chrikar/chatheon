package application

import (
	"errors"
	"testing"

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

type mockTokenGenerator struct{}

func (m *mockTokenGenerator) Generate(username, userID string) (string, error) {
	return "mock-token", nil
}

// ✅ Test: Happy path registration
func TestRegisterUser(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	err := service.Register("testuser", "password")
	assert.NoError(t, err)

	err = service.Register("testuser", "password")
	assert.ErrorIs(t, err, ErrUsernameTaken)
}

// ✅ Test: Empty username and password
func TestRegisterUser_EmptyFields(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	err := service.Register("", "password")
	assert.ErrorIs(t, err, ErrUsernameRequired)

	err = service.Register("username", "")
	assert.ErrorIs(t, err, ErrPasswordRequired)
}

// ✅ Test: Password is hashed (not plaintext)
func TestRegisterUser_PasswordIsHashed(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	username := "secureuser"
	password := "supersecure"
	err := service.Register(username, password)
	assert.NoError(t, err)

	storedUser, err := repo.FindByUsername(username)
	assert.NoError(t, err)
	assert.NotEqual(t, password, storedUser.PasswordHash, "Password should be hashed")
}

// ✅ Test: Register multiple users
func TestRegisterMultipleUsers(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

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

// ✅ Test: Successful login
func TestLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	err := service.Register("testuser", "password")
	assert.NoError(t, err)

	token, err := service.Login("testuser", "password")
	assert.NoError(t, err)
	assert.Equal(t, "mock-token", token)
}

// ✅ Test: Login failures
func TestLogin_Failure(t *testing.T) {
	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	err := service.Register("testuser", "password")
	assert.NoError(t, err)

	_, err = service.Login("unknown", "password")
	assert.ErrorIs(t, err, ErrInvalidCredentials)

	_, err = service.Login("testuser", "wrongpassword")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}
