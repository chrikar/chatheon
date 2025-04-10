package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/internal/auth"
)

type MessageHandler struct {
	messageService ports.MessageService
}

func NewMessageHandler(messageService ports.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

type createMessageRequest struct {
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
}

func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	senderID, ok := r.Context().Value(auth.ContextUserIDKey).(string)
	if !ok || senderID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Create message
	if err := h.messageService.CreateMessage(senderID, req.ReceiverID, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	limit := 10 // default
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed <= 0 {
			http.Error(w, "invalid 'limit' parameter: must be positive integer", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	offset := 0 // default
	if o := r.URL.Query().Get("offset"); o != "" {
		parsed, err := strconv.Atoi(o)
		if err != nil || parsed < 0 {
			http.Error(w, "invalid 'offset' parameter: must be non-negative integer", http.StatusBadRequest)
			return
		}
		offset = parsed
	}

	// Get user ID from context
	userID, ok := r.Context().Value(auth.ContextUserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch messages
	messages, err := h.messageService.GetMessagesByReceiver(userID, limit, offset)
	if err != nil {
		http.Error(w, "failed to fetch messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, "failed to encode messages", http.StatusInternalServerError)
		return
	}
}
