package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConversationFields(t *testing.T) {
	now := time.Now()
	id := uuid.New()
	conv := &Conversation{
		ID:             id,
		ParticipantIDs: []string{"alice", "bob"},
		CreatedAt:      now,
	}
	assert.Equal(t, id, conv.ID)
	assert.ElementsMatch(t, []string{"alice", "bob"}, conv.ParticipantIDs)
	assert.True(t, conv.CreatedAt.Equal(now))
}
