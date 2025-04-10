package ports

import "github.com/chrikar/chatheon/domain"

//go:generate mockery --name=MessageService --output=../../adapters/mocks --case=underscore
type MessageService interface {
	CreateMessage(senderID, receiverID, content string) error
	GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error)
}
