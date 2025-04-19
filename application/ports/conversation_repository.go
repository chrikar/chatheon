package ports

import "github.com/chrikar/chatheon/domain"

type ConversationRepository interface {
	Create(conversation *domain.Conversation) error
	FindByParticipant(userID string) ([]*domain.Conversation, error)
}
