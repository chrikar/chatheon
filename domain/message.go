package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID
	SenderID   string
	ReceiverID string
	Content    string
	CreatedAt  time.Time
}
