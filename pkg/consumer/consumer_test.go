package consumer

import (
	"context"
	"testing"
	"time"

	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/stretchr/testify/mock"
)

type MockNats struct {
	mock.Mock
}

func (m *MockNats) Publish(msg *infra.NatsMessage) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockNats) Subscribe(subject, consumerName string) error {
	return nil
}

func (m *MockNats) Fetch(messageCount int, ctx context.Context) ([]*infra.NatsMessage, error) {
	args := m.Called(messageCount, ctx)
	return args.Get(0).([]*infra.NatsMessage), args.Error(1)
}

func TestReturnSingleMessage(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	mockNats := new(MockNats)
	con := NewConsumer(mockNats)

	var messages []*infra.NatsMessage
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(`{\"TransactionID\":1,\"UserID\":1,\"Status\":\"OK\",\"Amount\": 456.89}`))

	messages = append(messages, msg)

	mockNats.On("Fetch", 100, ctx).Return(messages, nil)

	worker := make(chan bool)
	go con.Listen(ctx, worker, "TEST", "TEST_CONSUMER", 100)
	time.Sleep(10 * time.Millisecond)

	cancelFunc()
	<-worker

}
