// application/user_service_test.go
package application

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/internal/auth"
)

// mockUserRepo implements the repository port for Register().
type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) FindByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*domain.User), args.Error(1)
}
func (m *mockUserRepo) Create(user *domain.User) error {
	return m.Called(user).Error(0)
}

func TestUserService_Register(t *testing.T) {
	type scenario struct {
		name          string
		username      string
		password      string
		setupStubs    func(r *mockUserRepo)
		expectedError error
	}

	dbErr := errors.New("db failure")
	scenarios := []scenario{
		{
			name:          "empty username",
			username:      "",
			password:      "pw",
			setupStubs:    func(r *mockUserRepo) {},
			expectedError: ErrUsernameRequired,
		},
		{
			name:          "empty password",
			username:      "bob",
			password:      "",
			setupStubs:    func(r *mockUserRepo) {},
			expectedError: ErrPasswordRequired,
		},
		{
			name:     "duplicate username",
			username: "bob", password: "pw",
			setupStubs: func(r *mockUserRepo) {
				// stub FindByUsername to return an existing user
				r.On("FindByUsername", "bob").
					Return(&domain.User{ID: uuid.New(), Username: "bob"}, nil)
			},
			expectedError: ErrUsernameTaken,
		},
		{
			name:     "repo Create error",
			username: "carol", password: "pw",
			setupStubs: func(r *mockUserRepo) {
				r.On("FindByUsername", "carol").Return(nil, errors.New("user not found"))
				// capture the user passed in if you like
				r.On("Create", mock.MatchedBy(func(u *domain.User) bool {
					return u.Username == "carol"
				})).Return(dbErr)
			},
			expectedError: dbErr,
		},
		{
			name:     "success",
			username: "dave", password: "pw",
			setupStubs: func(r *mockUserRepo) {
				r.On("FindByUsername", "dave").Return(nil, errors.New("user not found"))
				r.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			repo := new(mockUserRepo)
			mgr := auth.NewJWTManager("secret", time.Hour)
			svc := NewUserService(repo, mgr)

			// arrange
			sc.setupStubs(repo)

			// act
			err := svc.Register(sc.username, sc.password)

			// assert
			if sc.expectedError != nil {
				assert.ErrorIs(t, err, sc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}
