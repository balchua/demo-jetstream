package infra

import (
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
)

type Nats interface {
	Publish(msg *nats.Msg) error
}

type JetStream interface {
	CreateStream(streamName string, streamOpts ...jsm.StreamOption) error
	CreateConsumer(streamName string, consumerOpts ...jsm.ConsumerOption) error
}
