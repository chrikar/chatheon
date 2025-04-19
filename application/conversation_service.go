package application

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

// ErrTooFewParticipants is returned when someone tries to create
// a conversation with fewer than two participants.
var ErrTooFewParticipants = errors.New("a conversation requires at least two participants")

// ConversationService is the application‑layer implementation
// of ports.ConversationService.
type ConversationService struct {
	repo ports.ConversationRepository
}

// NewConversationService constructs a ConversationService.
func NewConversationService(repo ports.ConversationRepository) *ConversationService {
	return &ConversationService{repo: repo}
}

// CreateConversation creates and persists a new conversation
// with the given participant IDs.
func (s *ConversationService) CreateConversation(participantIDs []string) (*domain.Conversation, error) {
	if len(participantIDs) < 2 {
		return nil, ErrTooFewParticipants
	}

	conv := &domain.Conversation{
		ID:             uuid.New(),
		ParticipantIDs: participantIDs,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.Create(conv); err != nil {
		return nil, err
	}
	return conv, nil
}

// GetConversationsForUser returns all conversations
// that include the given userID.
func (s *ConversationService) GetConversationsForUser(userID string) ([]*domain.Conversation, error) {
	return s.repo.FindByParticipant(userID)
}

// compile‑time check: ensure ConversationService implements the interface
var _ ports.ConversationService = (*ConversationService)(nil)
