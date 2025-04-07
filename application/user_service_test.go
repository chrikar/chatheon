package application

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/domain"
)

// Mock implementations

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

// ✅ Test: Register user (table-driven, thread-safe, parallelized)
func TestRegisterUser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		preRegister bool
		username   string
		password   string
		wantErr    error
	}{
		{"valid user", false, "testuser", "password", nil},
		{"duplicate username", true, "testuser", "password", ErrUsernameTaken},
		{"empty username", false, "", "password", ErrUsernameRequired},
		{"empty password", false, "testuser2", "", ErrPasswordRequired},
	}

	for _, tc := range cases {
		tc := tc // capture variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newMockUserRepo()
			service := NewUserService(repo, &mockTokenGenerator{})

			// Pre-register user for duplicate test
			if tc.preRegister {
				_ = service.Register(tc.username, tc.password)
			}

			err := service.Register(tc.username, tc.password)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ✅ Test: Password is hashed properly
func TestRegisterUser_PasswordIsHashed(t *testing.T) {
	t.Parallel()

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

// ✅ Test: Register and find multiple users safely
func TestRegisterMultipleUsers(t *testing.T) {
	t.Parallel()

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

	// Register users
	for _, u := range users {
		err := service.Register(u.username, u.password)
		assert.NoError(t, err)
	}

	// Verify users exist
	for _, u := range users {
		_, err := repo.FindByUsername(u.username)
		assert.NoError(t, err)
	}
}

// ✅ Test: Login flow (table-driven, parallel-safe)
func TestLogin(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		username   string
		password   string
		preRegister bool
		loginUser  string
		loginPass  string
		wantErr    error
	}{
		{"successful login", "testuser", "password", true, "testuser", "password", nil},
		{"wrong username", "testuser", "password", true, "unknown", "password", ErrInvalidCredentials},
		{"wrong password", "testuser", "password", true, "testuser", "wrongpassword", ErrInvalidCredentials},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := newMockUserRepo()
			service := NewUserService(repo, &mockTokenGenerator{})

			if tc.preRegister {
				err := service.Register(tc.username, tc.password)
				assert.NoError(t, err)
			}

			token, err := service.Login(tc.loginUser, tc.loginPass)

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "mock-token", token)
			}
		})
	}
}
