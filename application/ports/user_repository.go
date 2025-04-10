package ports

import "github.com/chrikar/chatheon/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByUsername(username string) (*domain.User, error)
}
