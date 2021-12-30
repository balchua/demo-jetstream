package infra

import (
	"context"

	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

type Nats interface {
	Publish(msg *NatsMessage) error
	Subscribe(subject, consumerName string) error
	Fetch(messageCount int, ctx context.Context) ([]*NatsMessage, error)
	Close()
}

type JetStream interface {
	CreateStream(streamName string, streamOpts ...jsm.StreamOption) error
	CreateConsumer(streamName string, consumerOpts ...jsm.ConsumerOption) error
	Close()
	IsStreamAvailable(streamName string) (bool, error)
}

type NatsMessage struct {
	msg *nats.Msg
}

func NewNatsMessage(subject string) *NatsMessage {

	msg := nats.NewMsg(subject)
	return &NatsMessage{
		msg: msg,
	}

}

func (m *NatsMessage) GetHeaders() map[string][]string {
	if m.msg == nil {
		return nil
	}
	return m.msg.Header
}

func (m *NatsMessage) GetHeader(key string) string {
	if m.msg == nil {
		return ""
	}
	return m.msg.Header.Get(key)
}

func (m *NatsMessage) GetBody() []byte {
	if m.msg == nil {
		return nil
	}
	return m.msg.Data
}

func (m *NatsMessage) AddHeader(key string, value string) {
	m.msg.Header.Add(key, value)
}

func (m *NatsMessage) SetBody(body []byte) {
	m.msg.Data = body
}

func (m *NatsMessage) Ack() {
	m.msg.Ack()
}

func (m *NatsMessage) Nack() {
	m.msg.Nak()
}

func (m *NatsMessage) GetUnderlyingNatsMessage() *nats.Msg {
	return m.msg
}
