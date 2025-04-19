package application

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/adapters/memory"
)

func TestConversationService_CreateAndList(t *testing.T) {
	repo := memory.NewConversationRepository()
	svc := NewConversationService(repo)

	// too few participants
	_, err := svc.CreateConversation([]string{"only-one"})
	assert.ErrorIs(t, err, ErrTooFewParticipants)

	// valid conversation
	ids := []string{"alice", "bob"}
	conv, err := svc.CreateConversation(ids)
	assert.NoError(t, err)
	assert.Equal(t, ids, conv.ParticipantIDs)
	assert.WithinDuration(t, time.Now(), conv.CreatedAt, time.Second)

	// list via service
	list, err := svc.GetConversationsForUser("alice")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, conv, list[0])

	// bob also sees it
	list, err = svc.GetConversationsForUser("bob")
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, conv, list[0])

	// user with no convs
	none, err := svc.GetConversationsForUser("charlie")
	assert.NoError(t, err)
	assert.Empty(t, none)
}
