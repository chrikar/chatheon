package application

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/ports"
)

var (
	ErrUsernameTaken = errors.New("username is already taken")
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(r ports.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(username, password string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}
	if password == "" {
		return errors.New("password cannot be empty")
	}

	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return ErrUsernameTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	return s.repo.Create(user)
}
