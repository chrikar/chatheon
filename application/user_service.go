package application

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

var (
	ErrUsernameTaken      = errors.New("username is already taken")
	ErrUsernameRequired   = errors.New("username cannot be empty")
	ErrPasswordRequired   = errors.New("password cannot be empty")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

type TokenGenerator interface {
	Generate(username, userID string) (string, error)
}

type UserService struct {
	repo     ports.UserRepository
	tokenGen TokenGenerator
}

func NewUserService(r ports.UserRepository, t TokenGenerator) *UserService {
	return &UserService{repo: r, tokenGen: t}
}

func (s *UserService) Register(username, password string) error {
	if username == "" {
		return ErrUsernameRequired
	}
	if password == "" {
		return ErrPasswordRequired
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

func (s *UserService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.tokenGen.Generate(user.Username, user.ID.String())
}
