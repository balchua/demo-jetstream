package infra

import (
	"context"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

type Nats interface {
	Publish(msg *nats.Msg) error
	Subscribe(subject, consumerName string) error
	Fetch(messageCount int, ctx context.Context) ([]*NatsMessage, error)
}

type JetStream interface {
	CreateStream(streamName string, streamOpts ...jsm.StreamOption) error
	CreateConsumer(streamName string, consumerOpts ...jsm.ConsumerOption) error
}

type NatsMessage struct {
	msg *nats.Msg
}

func NewNatsMessage(headers map[string][]string, body []byte) {

	//TODO

}

func (m *NatsMessage) GetHeaders() map[string][]string {
	return m.msg.Header
}

func (m *NatsMessage) GetHeader(key string) string {
	return m.msg.Header.Get(key)
}

func (m *NatsMessage) GetBody() []byte {
	return m.msg.Data
}

func (m *NatsMessage) Ack() {
	m.msg.Ack()
}

func (m *NatsMessage) Nack() {
	m.msg.Nak()
}
