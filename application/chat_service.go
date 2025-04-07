package application

import (
	"time"

	"github.com/google/uuid"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

type ChatService struct {
	repo     ports.MessageRepository
	notifier ports.NotificationService
}

func NewChatService(r ports.MessageRepository, n ports.NotificationService) *ChatService {
	return &ChatService{repo: r, notifier: n}
}

func (s *ChatService) SendMessage(fromUser, toUser, content string) error {
	msg := &domain.Message{
		ID:        uuid.New().String(),
		FromUser:  fromUser,
		ToUser:    toUser,
		Content:   content,
		Timestamp: time.Now(),
	}

	if err := s.repo.Save(msg); err != nil {
		return err
	}

	return s.notifier.Notify(toUser, "You have a new message!")
}

func (s *ChatService) GetMessages(userID string) ([]*domain.Message, error) {
	return s.repo.GetMessagesForUser(userID)
}