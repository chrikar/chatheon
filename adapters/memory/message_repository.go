package memory

import (
	"sync"

	"github.com/chrikar/chatheon/domain"
)

type MessageRepository struct {
	mu       sync.RWMutex
	messages []*domain.Message
}

func NewMessageRepository() *MessageRepository {
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

func (r *MessageRepository) GetMessagesBySender(senderID string) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Message
	for _, msg := range r.messages {
		if msg.SenderID == senderID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (r *MessageRepository) GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Message
	for _, msg := range r.messages {
		if msg.ReceiverID == receiverID {
			result = append(result, msg)
		}
	}

	// Apply pagination
	start := offset
	if start > len(result) {
		start = len(result)
	}
	end := start + limit
	if end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}
