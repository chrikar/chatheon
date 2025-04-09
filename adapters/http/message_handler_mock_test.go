package http

import (
	"github.com/stretchr/testify/mock"

	"github.com/chrikar/chatheon/application/ports"
	"github.com/chrikar/chatheon/domain"
)

type mockMessageService struct {
	mock.Mock
}

var _ ports.MessageService = (*mockMessageService)(nil) // Optional compile-time check

func (m *mockMessageService) CreateMessage(senderID, receiverID, content string) error {
	args := m.Called(senderID, receiverID, content)
	return args.Error(0)
}

func (m *mockMessageService) GetMessagesByReceiver(receiverID string) ([]*domain.Message, error) {
	args := m.Called(receiverID)
	return args.Get(0).([]*domain.Message), args.Error(1)
}
