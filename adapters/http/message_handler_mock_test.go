package http

import (
	"github.com/stretchr/testify/mock"

	"github.com/chrikar/chatheon/domain"
)

type mockMessageService struct {
	mock.Mock
}

func (m *mockMessageService) CreateMessage(senderID, receiverID, content string) error {
	args := m.Called(senderID, receiverID, content)
	return args.Error(0)
}

func (m *mockMessageService) GetMessagesByReceiver(receiverID string) ([]*domain.Message, error) {
	args := m.Called(receiverID)
	return args.Get(0).([]*domain.Message), args.Error(1)
}
