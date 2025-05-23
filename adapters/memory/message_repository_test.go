package memory

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/domain"
)

func TestMessageRepository_CreateAndGetMessagesBySender(t *testing.T) {
	t.Parallel()

	repo := NewMessageRepository()

	// Prepare messages
	message1 := &domain.Message{
		ID:         uuid.New(),
		SenderID:   "user-1",
		ReceiverID: "user-2",
		Content:    "Hello, user-2!",
	}
	message2 := &domain.Message{
		ID:         uuid.New(),
		SenderID:   "user-1",
		ReceiverID: "user-3",
		Content:    "Hello, user-3!",
	}
	message3 := &domain.Message{
		ID:         uuid.New(),
		SenderID:   "user-2",
		ReceiverID: "user-1",
		Content:    "Hi, user-1!",
	}

	// Create messages
	err := repo.Create(message1)
	assert.NoError(t, err)

	err = repo.Create(message2)
	assert.NoError(t, err)

	err = repo.Create(message3)
	assert.NoError(t, err)

	// Get messages by sender
	user1Messages, err := repo.GetMessagesBySender("user-1")
	assert.NoError(t, err)
	assert.Len(t, user1Messages, 2)

	user2Messages, err := repo.GetMessagesBySender("user-2")
	assert.NoError(t, err)
	assert.Len(t, user2Messages, 1)

	unknownUserMessages, err := repo.GetMessagesBySender("unknown")
	assert.NoError(t, err)
	assert.Len(t, unknownUserMessages, 0)
}

func TestMessageRepository_GetMessagesByReceiver(t *testing.T) {
	t.Parallel()

	repo := NewMessageRepository()

	// Prepare messages
	err := repo.Create(&domain.Message{ID: uuid.New(), SenderID: "user-1", ReceiverID: "user-2", Content: "Hi user2!"})
	assert.NoError(t, err)
	err = repo.Create(&domain.Message{ID: uuid.New(), SenderID: "user-3", ReceiverID: "user-2", Content: "Hello user2!"})
	assert.NoError(t, err)
	err = repo.Create(&domain.Message{ID: uuid.New(), SenderID: "user-1", ReceiverID: "user-3", Content: "Hi user3!"})
	assert.NoError(t, err)

	messages, err := repo.GetMessagesByReceiver("user-2", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	messages, err = repo.GetMessagesByReceiver("user-3", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)

	messages, err = repo.GetMessagesByReceiver("user-unknown", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, messages, 0)
}

func TestMessageRepository_Pagination(t *testing.T) {
	repo := NewMessageRepository()

	receiverID := "receiver-1"
	for i := 0; i < 10; i++ {
		_ = repo.Create(&domain.Message{
			ID:         uuid.New(),
			SenderID:   fmt.Sprintf("sender-%d", i),
			ReceiverID: receiverID,
			Content:    fmt.Sprintf("Message %d", i),
		})
	}

	tests := []struct {
		limit, offset int
		expectedCount int
	}{
		{limit: 5, offset: 0, expectedCount: 5},
		{limit: 3, offset: 7, expectedCount: 3},
		{limit: 5, offset: 10, expectedCount: 0},
	}

	for _, tc := range tests {
		messages, err := repo.GetMessagesByReceiver(receiverID, tc.limit, tc.offset)
		assert.NoError(t, err)
		assert.Len(t, messages, tc.expectedCount)
	}
}
