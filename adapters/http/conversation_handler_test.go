package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chrikar/chatheon/domain"
)

// mockConversationService is a testify/mock for ports.ConversationService.
type mockConversationService struct {
	mock.Mock
}

func (m *mockConversationService) CreateConversation(participantIDs []string) (*domain.Conversation, error) {
	args := m.Called(participantIDs)
	return args.Get(0).(*domain.Conversation), args.Error(1)
}

func (m *mockConversationService) GetConversationsForUser(userID string) ([]*domain.Conversation, error) {
	args := m.Called(userID)
	return args.Get(0).([]*domain.Conversation), args.Error(1)
}

func TestConversationHandler_CreateConversation_Success(t *testing.T) {
	service := new(mockConversationService)
	handler := NewConversationHandler(service)

	// prepare expected conversation
	now := time.Now()
	ids := []string{"alice", "bob"}
	conv := &domain.Conversation{
		ID:             uuid.New(),
		ParticipantIDs: ids,
		CreatedAt:      now,
	}
	service.On("CreateConversation", ids).Return(conv, nil)

	// build request with both participants
	reqBody := createConversationRequest{ParticipantIDs: ids}
	payload, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/conversations", bytes.NewReader(payload))
	req = req.WithContext(contextWithUserID(req.Context(), "alice"))
	rr := httptest.NewRecorder()

	handler.CreateConversation(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var got domain.Conversation
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, conv.ID, got.ID)
	assert.ElementsMatch(t, conv.ParticipantIDs, got.ParticipantIDs)
	assert.WithinDuration(t, conv.CreatedAt, got.CreatedAt, time.Second)

	service.AssertExpectations(t)
}

func TestConversationHandler_CreateConversation_TooFewParticipants(t *testing.T) {
	service := new(mockConversationService)
	handler := NewConversationHandler(service)

	// only one participant in request → handler should reject first
	reqBody := createConversationRequest{ParticipantIDs: []string{"alice"}}
	payload, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/conversations", bytes.NewReader(payload))
	req = req.WithContext(contextWithUserID(req.Context(), "alice"))
	rr := httptest.NewRecorder()

	handler.CreateConversation(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	service.AssertNotCalled(t, "CreateConversation", mock.Anything)
}

func TestConversationHandler_GetConversations_Success(t *testing.T) {
	service := new(mockConversationService)
	handler := NewConversationHandler(service)

	now := time.Now()
	expected := []*domain.Conversation{
		{ID: uuid.New(), ParticipantIDs: []string{"alice", "bob"}, CreatedAt: now},
	}
	service.On("GetConversationsForUser", "alice").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/conversations", nil)
	req = req.WithContext(contextWithUserID(req.Context(), "alice"))
	rr := httptest.NewRecorder()

	handler.GetConversations(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var got []*domain.Conversation
	assert.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Len(t, got, 1)
	assert.Equal(t, expected[0].ID, got[0].ID)
	assert.ElementsMatch(t, expected[0].ParticipantIDs, got[0].ParticipantIDs)
	assert.WithinDuration(t, expected[0].CreatedAt, got[0].CreatedAt, time.Second)

	service.AssertExpectations(t)
}
