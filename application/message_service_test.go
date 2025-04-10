package application

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/domain"
)

// Mock message repository
type mockMessageRepo struct {
	messages []*domain.Message
}

func newMockMessageRepo() *mockMessageRepo {
	return &mockMessageRepo{messages: []*domain.Message{}}
}

func (m *mockMessageRepo) Create(message *domain.Message) error {
	if message.Content == "fail" {
		return errors.New("forced failure")
	}
	m.messages = append(m.messages, message)
	return nil
}

func (m *mockMessageRepo) GetMessagesBySender(senderID string) ([]*domain.Message, error) {
	var result []*domain.Message
	for _, msg := range m.messages {
		if msg.SenderID == senderID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (m *mockMessageRepo) GetMessagesByReceiver(receiverID string) ([]*domain.Message, error) {
	var result []*domain.Message
	for _, msg := range m.messages {
		if msg.ReceiverID == receiverID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func TestCreateMessage(t *testing.T) {
	t.Parallel()

	repo := newMockMessageRepo()
	service := NewMessageService(repo)

	err := service.CreateMessage("user-1", "user-2", "Hello, world!")
	assert.NoError(t, err)
	assert.Len(t, repo.messages, 1)

	err = service.CreateMessage("user-1", "user-2", "")
	assert.ErrorIs(t, err, ErrMessageContentRequired)

	err = service.CreateMessage("user-1", "user-2", "fail")
	assert.Error(t, err)
}

func TestGetMessages(t *testing.T) {
	t.Parallel()

	repo := newMockMessageRepo()
	service := NewMessageService(repo)

	// Pre-fill messages
	err := service.CreateMessage("user-1", "user-2", "Hello!")
	assert.NoError(t, err)
	err = service.CreateMessage("user-2", "user-1", "Hi!")
	assert.NoError(t, err)
	err = service.CreateMessage("user-1", "user-2", "How are you?")
	assert.NoError(t, err)

	messages, err := service.GetMessages("user-1")
	assert.NoError(t, err)
	assert.Len(t, messages, 2)

	messages, err = service.GetMessages("user-2")
	assert.NoError(t, err)
	assert.Len(t, messages, 1)

	messages, err = service.GetMessages("nonexistent")
	assert.NoError(t, err)
	assert.Len(t, messages, 0)
}
