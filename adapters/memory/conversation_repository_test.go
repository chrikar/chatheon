package memory

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/domain"
)

func TestConversationRepository_CreateAndFindByParticipant(t *testing.T) {
	repo := NewConversationRepository()

	// prepare three conversations
	now := time.Now()
	conv1 := &domain.Conversation{
		ID:             uuid.New(),
		ParticipantIDs: []string{"alice", "bob"},
		CreatedAt:      now,
	}
	conv2 := &domain.Conversation{
		ID:             uuid.New(),
		ParticipantIDs: []string{"bob", "carol"},
		CreatedAt:      now,
	}
	conv3 := &domain.Conversation{
		ID:             uuid.New(),
		ParticipantIDs: []string{"dave", "eve"},
		CreatedAt:      now,
	}

	// create them
	assert.NoError(t, repo.Create(conv1))
	assert.NoError(t, repo.Create(conv2))
	assert.NoError(t, repo.Create(conv3))

	// bob participates in conv1 & conv2
	bobConvs, err := repo.FindByParticipant("bob")
	assert.NoError(t, err)
	assert.Len(t, bobConvs, 2)
	assert.Contains(t, bobConvs, conv1)
	assert.Contains(t, bobConvs, conv2)

	// alice only in conv1
	aliceConvs, err := repo.FindByParticipant("alice")
	assert.NoError(t, err)
	assert.Len(t, aliceConvs, 1)
	assert.Equal(t, conv1, aliceConvs[0])

	// unknown user gets none
	unk, err := repo.FindByParticipant("unknown")
	assert.NoError(t, err)
	assert.Empty(t, unk)
}
