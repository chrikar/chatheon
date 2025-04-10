package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/chrikar/chatheon/adapters/mocks"
	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/internal/auth"
)

// Test helper for context injection
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, auth.ContextUserIDKey, userID)
}

func TestMessageHandler_CreateMessage(t *testing.T) {
	service := new(mocks.MessageService)
	service.On("CreateMessage", "user-1", "user-2", "Hello!").Return(nil)

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

	service.AssertExpectations(t)
}

func TestMessageHandler_GetMessages(t *testing.T) {
	service := new(mocks.MessageService)
	expectedMessages := []*domain.Message{
		{SenderID: "user-2", ReceiverID: "user-1", Content: "Hi user1!", CreatedAt: time.Now()},
		{SenderID: "user-3", ReceiverID: "user-1", Content: "Hello user1!", CreatedAt: time.Now()},
	}
	service.On("GetMessagesByReceiver", "user-1", 10, 0).Return(expectedMessages, nil)

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
	assert.NotEmpty(t, response[0].CreatedAt)
	assert.NotEmpty(t, response[1].CreatedAt)

	for _, msg := range response {
		_, err := uuid.Parse(msg.ID.String())
		assert.NoError(t, err, "ID should be valid UUID")
	}

	service.AssertExpectations(t)
}

func TestMessageHandler_GetMessages_Pagination(t *testing.T) {
	service := new(mocks.MessageService)
	expectedMessages := []*domain.Message{
		{ID: uuid.New(), SenderID: "sender-1", ReceiverID: "user-1", Content: "Hi 1", CreatedAt: time.Now()},
		{ID: uuid.New(), SenderID: "sender-2", ReceiverID: "user-1", Content: "Hi 2", CreatedAt: time.Now()},
	}
	service.On("GetMessagesByReceiver", "user-1", 2, 0).Return(expectedMessages, nil)

	handler := NewMessageHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/messages?limit=2&offset=0", nil)
	req = req.WithContext(contextWithUserID(req.Context(), "user-1"))
	rr := httptest.NewRecorder()

	handler.GetMessages(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []domain.Message
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)

	service.AssertExpectations(t)
}

func TestMessageHandler_GetMessages_InvalidPagination(t *testing.T) {
	service := new(mocks.MessageService)
	handler := NewMessageHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/messages?limit=abc&offset=-1", nil)
	req = req.WithContext(contextWithUserID(req.Context(), "user-1"))
	rr := httptest.NewRecorder()

	handler.GetMessages(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
