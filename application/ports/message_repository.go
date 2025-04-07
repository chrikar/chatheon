package ports

import "github.com/chrikar/chatheon/domain"

type MessageRepository interface {
	Save(message *domain.Message) error
	GetMessagesForUser(userID string) ([]*domain.Message, error)
}