package memory

import (
	"sync"

	"github.com/chrikar/chatheon/domain"
)

// ConversationRepository is an in‑memory implementation of ports.ConversationRepository.
type ConversationRepository struct {
	mu            sync.RWMutex
	conversations []*domain.Conversation
}

// NewConversationRepository constructs an in‑memory repo.
func NewConversationRepository() *ConversationRepository {
	return &ConversationRepository{
		conversations: make([]*domain.Conversation, 0),
	}
}

// Create appends a new conversation.
func (r *ConversationRepository) Create(conv *domain.Conversation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conversations = append(r.conversations, conv)
	return nil
}

// FindByParticipant filters conversations by userID.
func (r *ConversationRepository) FindByParticipant(userID string) ([]*domain.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Conversation
	for _, c := range r.conversations {
		for _, pid := range c.ParticipantIDs {
			if pid == userID {
				result = append(result, c)
				break
			}
		}
	}
	return result, nil
}
