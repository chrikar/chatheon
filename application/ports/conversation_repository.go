package ports

import "github.com/chrikar/chatheon/domain"

// ConversationRepository defines persistence for conversations.
type ConversationRepository interface {
	// Create persists a new conversation.
	Create(conversation *domain.Conversation) error
	// FindByParticipant returns all conversations containing userID.
	FindByParticipant(userID string) ([]*domain.Conversation, error)
}
