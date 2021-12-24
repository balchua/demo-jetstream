package infra

import "github.com/nats-io/nats.go"

type NatsInfo struct {
	Nc *nats.Conn
	Js nats.JetStreamContext
}
