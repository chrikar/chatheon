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

// ✅ Table-driven test: Registration
func TestRegisterUser(t *testing.T) {
	t.Parallel()

	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	cases := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{"valid user", "testuser", "password", nil},
		{"duplicate username", "testuser", "password", ErrUsernameTaken},
		{"empty username", "", "password", ErrUsernameRequired},
		{"empty password", "testuser2", "", ErrPasswordRequired},
	}

	// First register to setup duplicate test
	_ = service.Register("testuser", "password")

	for _, tc := range cases {
		tc := tc // capture variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a fresh repo & service for each test case
			repo := newMockUserRepo()
			service := NewUserService(repo, &mockTokenGenerator{})

			// For "duplicate username" case, pre-register
			if tc.wantErr == ErrUsernameTaken {
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

// ✅ Test: Password is hashed (explicit check)
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

// ✅ Test: Register multiple users
func TestRegisterMultipleUsers(t *testing.T) {
	t.Parallel()

	// Each subtest gets its own fresh state
	t.Run("register and find multiple users", func(t *testing.T) {
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
	})
}

// ✅ Table-driven test: Login
func TestLogin(t *testing.T) {
	t.Parallel()

	repo := newMockUserRepo()
	service := NewUserService(repo, &mockTokenGenerator{})

	// Pre-register user
	err := service.Register("testuser", "password")
	assert.NoError(t, err)

	cases := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{"successful login", "testuser", "password", nil},
		{"wrong username", "unknown", "password", ErrInvalidCredentials},
		{"wrong password", "testuser", "wrongpassword", ErrInvalidCredentials},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			token, err := service.Login(tc.username, tc.password)

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
