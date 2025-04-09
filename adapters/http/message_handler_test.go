package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/internal/auth"
)

// Mock service for testing
type mockMessageService struct {
	messages []*domain.Message
}

func (m *mockMessageService) CreateMessage(senderID, receiverID, content string) error {
	m.messages = append(m.messages, &domain.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	})
	return nil
}

func (m *mockMessageService) GetMessagesByReceiver(receiverID string) ([]*domain.Message, error) {
	var result []*domain.Message
	for _, msg := range m.messages {
		if msg.ReceiverID == receiverID {
			result = append(result, msg)
		}
	}
	return result, nil
}

// Test helper for context injection
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, auth.ContextUserIDKey, userID)
}

func TestMessageHandler_CreateMessage(t *testing.T) {
	service := &mockMessageService{}
	handler := NewMessageHandler(service)

	reqBody := createMessageRequest{
		ReceiverID: "user-2",
		Content:    "Hello!",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewReader(body))
	ctx := contextWithUserID(req.Context(), "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.CreateMessage(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Len(t, service.messages, 1)
	assert.Equal(t, "user-1", service.messages[0].SenderID)
	assert.Equal(t, "user-2", service.messages[0].ReceiverID)
	assert.Equal(t, "Hello!", service.messages[0].Content)
}

func TestMessageHandler_GetMessages(t *testing.T) {
	service := &mockMessageService{
		messages: []*domain.Message{
			{SenderID: "user-2", ReceiverID: "user-1", Content: "Hi user1!"},
			{SenderID: "user-3", ReceiverID: "user-1", Content: "Hello user1!"},
			{SenderID: "user-1", ReceiverID: "user-3", Content: "Hi user3!"},
		},
	}
	handler := NewMessageHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	ctx := contextWithUserID(req.Context(), "user-1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.GetMessages(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []*domain.Message
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Hi user1!", response[0].Content)
	assert.Equal(t, "Hello user1!", response[1].Content)
}
