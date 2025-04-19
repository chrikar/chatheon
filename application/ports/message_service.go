package ports

import "github.com/chrikar/chatheon/domain"

type MessageService interface {
	CreateMessage(senderID, receiverID, content string) error
	GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error)
	SetMessageStatus(messageID string, status domain.MessageStatus) error
}
