package memory

import (
	"sync"

	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/ports"
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

func (r *MessageRepository) Save(message *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, message)
	return nil
}

func (r *MessageRepository) GetMessagesForUser(userID string) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var userMessages []*domain.Message
	for _, msg := range r.messages {
		if msg.ToUser == userID || msg.FromUser == userID {
			userMessages = append(userMessages, msg)
		}
	}
	return userMessages, nil
}