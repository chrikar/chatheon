// application/message_service_test.go
package application

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chrikar/chatheon/domain"
)

// mockMessageRepo implements the ports.MessageRepository interface.
type mockMessageRepo struct{ mock.Mock }

func (m *mockMessageRepo) Create(msg *domain.Message) error {
	return m.Called(msg).Error(0)
}

func (m *mockMessageRepo) GetMessagesByReceiver(receiverID string, limit, offset int) ([]*domain.Message, error) {
	args := m.Called(receiverID, limit, offset)
	return args.Get(0).([]*domain.Message), args.Error(1)
}

// satisfy the interface: this method exists but isn't used by our service directly
func (m *mockMessageRepo) GetMessagesBySender(senderID string) ([]*domain.Message, error) {
	args := m.Called(senderID)
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *mockMessageRepo) SetMessageStatus(id uuid.UUID, status domain.MessageStatus) error {
	return m.Called(id, status).Error(0)
}

func TestMessageService_CreateMessage(t *testing.T) {
	dbFailError := errors.New("db fail")
	tests := []struct {
		name          string
		sender        string
		receiver      string
		content       string
		setupStubs    func(*mockMessageRepo)
		expectedError error
	}{
		{
			name:          "empty content",
			sender:        "u1",
			receiver:      "u2",
			content:       "",
			setupStubs:    func(r *mockMessageRepo) {},
			expectedError: ErrMessageContentRequired,
		},
		{
			name:     "repo error",
			sender:   "u1",
			receiver: "u2",
			content:  "hello",
			setupStubs: func(r *mockMessageRepo) {
				r.On("Create", mock.AnythingOfType("*domain.Message")).
					Return(dbFailError)
			},
			expectedError: dbFailError,
		},
		{
			name:     "success",
			sender:   "u1",
			receiver: "u2",
			content:  "hi!",
			setupStubs: func(r *mockMessageRepo) {
				r.On("Create", mock.MatchedBy(func(m *domain.Message) bool {
					return m.SenderID == "u1" &&
						m.ReceiverID == "u2" &&
						m.Content == "hi!" &&
						m.ID != uuid.Nil
				})).Return(nil)
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mockMessageRepo)
			svc := NewMessageService(repo)

			tc.setupStubs(repo)
			err := svc.CreateMessage(tc.sender, tc.receiver, tc.content)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
		})
	}
}

func TestMessageService_GetMessagesByReceiver(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := NewMessageService(repo)

	now := time.Now()
	fake := []*domain.Message{
		{ID: uuid.New(), SenderID: "u1", ReceiverID: "u2", Content: "a", CreatedAt: now, Status: domain.StatusSent},
	}
	repo.On("GetMessagesByReceiver", "u2", 5, 1).Return(fake, nil)

	out, err := svc.GetMessagesByReceiver("u2", 5, 1)
	assert.NoError(t, err)
	assert.Equal(t, fake, out)

	repo.AssertExpectations(t)
}

func TestMessageService_SetMessageStatus(t *testing.T) {
	repo := new(mockMessageRepo)
	svc := NewMessageService(repo)

	// invalid UUID
	errInvalid := svc.SetMessageStatus("not-uuid", domain.StatusRead)
	assert.Error(t, errInvalid)

	// repo error
	id := uuid.New()
	dbErr := errors.New("not found")
	repo.On("SetMessageStatus", id, domain.StatusDelivered).Return(dbErr)
	err := svc.SetMessageStatus(id.String(), domain.StatusDelivered)
	assert.ErrorIs(t, err, dbErr)
	repo.AssertExpectations(t)

	// success
	repo.ExpectedCalls = nil
	repo.On("SetMessageStatus", id, domain.StatusRead).Return(nil)
	errSuccess := svc.SetMessageStatus(id.String(), domain.StatusRead)
	assert.NoError(t, errSuccess)
}
