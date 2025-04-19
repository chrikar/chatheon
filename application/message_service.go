package application

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

type MessageServiceInterface interface {
	CreateMessage(senderID, receiverID, content string) error
	GetMessagesByReceiver(receiverID string) ([]*domain.Message, error)
	SetMessageStatus(messageID string, status domain.MessageStatus) error
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
	}

	return s.repo.Create(message)
}

func (s *MessageService) GetMessages(senderID string) ([]*domain.Message, error) {
	return s.repo.GetMessagesBySender(senderID)
}

func (s *MessageService) GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error) {
	return s.repo.GetMessagesByReceiver(receiverID, limit, offset)
}

func (s *MessageService) SetMessageStatus(messageID string, status domain.MessageStatus) error {
	id, err := uuid.Parse(messageID)
	if err != nil {
		return fmt.Errorf("invalid message ID: %w", err)
	}
	return s.repo.SetMessageStatus(id, status)
}

var _ ports.MessageService = (*MessageService)(nil)
