package ports

import "github.com/chrikar/chatheon/domain"

type MessageService interface {
	CreateMessage(senderID, receiverID, content string) error
	GetMessagesByReceiver(receiverID string) ([]*domain.Message, error)
}
