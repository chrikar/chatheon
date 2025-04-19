package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
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
	senderID, ok := r.Context().Value(auth.ContextUserIDKey).(string)
	if !ok || senderID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req createMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.messageService.CreateMessage(senderID, req.ReceiverID, req.Content); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *MessageHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var payload struct {
		Status domain.MessageStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.messageService.SetMessageStatus(id, payload.Status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.ContextUserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// default pagination
	limit := 10
	offset := 0

	// parse limit
	if l := r.URL.Query().Get("limit"); l != "" {
		n, err := strconv.Atoi(l)
		if err != nil || n <= 0 {
			http.Error(w,
				"invalid 'limit' parameter: must be a positive integer",
				http.StatusBadRequest,
			)
			return
		}
		limit = n
	}

	// parse offset
	if o := r.URL.Query().Get("offset"); o != "" {
		n, err := strconv.Atoi(o)
		if err != nil || n < 0 {
			http.Error(w,
				"invalid 'offset' parameter: must be a non-negative integer",
				http.StatusBadRequest,
			)
			return
		}
		offset = n
	}

	msgs, err := h.messageService.GetMessagesByReceiver(userID, limit, offset)
	if err != nil {
		http.Error(w, "failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// build JSON response
	type messageResponse struct {
		ID         string `json:"id"`
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
		CreatedAt  string `json:"created_at"`
		Status     string `json:"status"`
	}

	resp := make([]messageResponse, 0, len(msgs))
	for _, m := range msgs {
		resp = append(resp, messageResponse{
			ID:         m.ID.String(),
			SenderID:   m.SenderID,
			ReceiverID: m.ReceiverID,
			Content:    m.Content,
			CreatedAt:  m.CreatedAt.Format(time.RFC3339Nano),
			Status:     m.Status.String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
