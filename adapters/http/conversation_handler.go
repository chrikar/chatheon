package http

import (
	"encoding/json"
	"net/http"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/internal/auth"
)

type ConversationHandler struct {
	svc ports.ConversationService
}

func NewConversationHandler(svc ports.ConversationService) *ConversationHandler {
	return &ConversationHandler{svc: svc}
}

type createConversationRequest struct {
	ParticipantIDs []string `json:"participant_ids"`
}

func (h *ConversationHandler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	// Auth
	userID, ok := r.Context().Value(auth.ContextUserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Decode
	var req createConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Ensure current user included
	if !contains(req.ParticipantIDs, userID) {
		req.ParticipantIDs = append(req.ParticipantIDs, userID)
	}

	// Must have at least 2 participants
	if len(req.ParticipantIDs) < 2 {
		http.Error(w, "a conversation requires at least two participants", http.StatusBadRequest)
		return
	}

	// Call service
	conv, err := h.svc.CreateConversation(req.ParticipantIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(conv)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *ConversationHandler) GetConversations(w http.ResponseWriter, r *http.Request) {
	// Auth
	userID, ok := r.Context().Value(auth.ContextUserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Call service
	convs, err := h.svc.GetConversationsForUser(userID)
	if err != nil {
		http.Error(w, "failed to fetch conversations", http.StatusInternalServerError)
		return
	}

	// Respond
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(convs)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// Helper
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
