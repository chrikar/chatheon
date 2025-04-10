package application

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

type MessageServiceInterface interface {
	CreateMessage(senderID, receiverID, content string) error
	GetMessagesByReceiver(receiverID string) ([]*domain.Message, error)
}

var (
	ErrMessageContentRequired = errors.New("message content cannot be empty")
)

type MessageService struct {
	repo ports.MessageRepository
}

func NewMessageService(repo ports.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) CreateMessage(senderID, receiverID, content string) error {
	if content == "" {
		return ErrMessageContentRequired
	}

	message := &domain.Message{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	return s.repo.Create(message)
}

func (s *MessageService) GetMessages(senderID string) ([]*domain.Message, error) {
	return s.repo.GetMessagesBySender(senderID)
}

func (s *MessageService) GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error) {
	return s.repo.GetMessagesByReceiver(receiverID, limit, offset)
}
