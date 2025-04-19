package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chrikar/chatheon/adapters/mocks"
	"github.com/chrikar/chatheon/domain"
	"github.com/chrikar/chatheon/internal/auth"
)

// mockMessageService is the testify/mock generated or hand‐rolled mock.
type mockMessageService struct {
	mock.Mock
}

func (m *mockMessageService) CreateMessage(senderID, receiverID, content string) error {
	args := m.Called(senderID, receiverID, content)
	return args.Error(0)
}

func (m *mockMessageService) GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error) {
	args := m.Called(receiverID, limit, offset)
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *mockMessageService) SetMessageStatus(messageID string, status domain.MessageStatus) error {
	args := m.Called(messageID, status)
	return args.Error(0)
}

// helper to inject user ID into request context
func contextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, auth.ContextUserIDKey, userID)
}

func TestMessageHandler_GetMessages_InvalidPagination(t *testing.T) {
	service := new(mockMessageService)
	handler := NewMessageHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/messages?limit=abc&offset=-1", nil)
	req = req.WithContext(contextWithUserID(req.Context(), "user-1"))
	rr := httptest.NewRecorder()

	handler.GetMessages(rr, req)

	// Should reject the bad params with 400 and never call the service
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	service.AssertNotCalled(t, "GetMessagesByReceiver", mock.Anything, mock.Anything, mock.Anything)
}

func TestMessageHandler_GetMessages_Success(t *testing.T) {
	now := time.Now()
	expected := []*domain.Message{
		{ID: uuid.New(), SenderID: "s1", ReceiverID: "user-1", Content: "hi1", CreatedAt: now, Status: domain.StatusDelivered},
	}

	service := mocks.NewMockMessageService(t)
	service.On("GetMessagesByReceiver", "user-1", 10, 0).Return(expected, nil)

	handler := NewMessageHandler(service)
	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	req = req.WithContext(contextWithUserID(req.Context(), "user-1"))
	rr := httptest.NewRecorder()

	handler.GetMessages(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var got []domain.Message
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Len(t, got, 1)
	assert.Equal(t, expected[0].Status, got[0].Status) // enum field round‑tripped

	service.AssertExpectations(t)
}
