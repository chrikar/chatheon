package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MessageStatus is an iota enum for message states.
type MessageStatus int

const (
	StatusSent MessageStatus = iota
	StatusDelivered
	StatusRead
)

var statusNames = []string{"sent", "delivered", "read"}

// String implements fmt.Stringer.
func (s MessageStatus) String() string {
	if s < 0 || int(s) >= len(statusNames) {
		return "unknown"
	}
	return statusNames[s]
}

// MarshalJSON renders the enum as a JSON string.
func (s MessageStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON parses a JSON string into the enum.
func (s *MessageStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	for i, name := range statusNames {
		if name == str {
			*s = MessageStatus(i)
			return nil
		}
	}
	return fmt.Errorf("invalid MessageStatus %q", str)
}

type Message struct {
	ID         uuid.UUID     `json:"id"`
	SenderID   string        `json:"sender_id"`
	ReceiverID string        `json:"receiver_id"`
	Content    string        `json:"content"`
	CreatedAt  time.Time     `json:"created_at"`
	Status     MessageStatus `json:"status"`
}
