package domain

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID             uuid.UUID `json:"id"`
	ParticipantIDs []string  `json:"participant_ids"`
	CreatedAt      time.Time `json:"created_at"`
}
