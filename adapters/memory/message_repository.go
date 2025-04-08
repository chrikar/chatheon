package memory

import (
	"sync"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

type MessageRepository struct {
	messages []*domain.Message
	mu       sync.RWMutex
}

func NewMessageRepository() ports.MessageRepository {
	return &MessageRepository{
		messages: make([]*domain.Message, 0),
	}
}

func (r *MessageRepository) Create(message *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, message)
	return nil
}

func (r *MessageRepository) GetMessagesBySender(userID string) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMessages []*domain.Message
	for _, msg := range r.messages {
		if msg.SenderID == userID {
			userMessages = append(userMessages, msg)
		}
	}
	return userMessages, nil
}