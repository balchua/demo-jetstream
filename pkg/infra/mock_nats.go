package infra

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockNats struct {
	mock.Mock
}

func (m *MockNats) Publish(msg *NatsMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockNats) Subscribe(subject, consumerName string) error {
	args := m.Called(subject, consumerName)
	return args.Error(0)
}

func (m *MockNats) Fetch(messageCount int, ctx context.Context) ([]*NatsMessage, error) {
	args := m.Called(messageCount, ctx)
	return args.Get(0).([]*NatsMessage), args.Error(1)
}

func (m *MockNats) Close() {
	return
}
