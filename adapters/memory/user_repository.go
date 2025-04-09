package memory

import (
	"errors"
	"sync"

	"github.com/chrikar/chatheon/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return errors.New("user already exists")
	}

	r.users[user.Username] = user
	return nil
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
