package ports

import "github.com/chrikar/chatheon/domain"

// ConversationService handles creating and listing conversations.
type ConversationService interface {
	// Create a new conversation with the given participants.
	CreateConversation(participantIDs []string) (*domain.Conversation, error)

	// List all conversations that a user participates in.
	GetConversationsForUser(userID string) ([]*domain.Conversation, error)
}
