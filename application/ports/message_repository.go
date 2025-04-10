package ports

import "github.com/chrikar/chatheon/domain"

type MessageRepository interface {
	Create(message *domain.Message) error
	GetMessagesBySender(senderID string) ([]*domain.Message, error)
	GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error)
}
