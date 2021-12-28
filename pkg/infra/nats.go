package infra

import (
	"errors"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NatsInfo struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewNats(seedFile string, natsUri string) (*NatsInfo, error) {
	opt, err := nats.NkeyOptionFromSeed("hack/seed.txt")

	if err != nil {
		zap.S().Fatalf("unable to get nkey seed %v", err)
	}

	nc, err := nats.Connect("localhost:32422", opt)
	if err != nil {
		zap.S().Fatalf("unable to connect to nats server %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		zap.S().Fatalf("unable to get jetstream context %v", err)
	}

	natsInfo := &NatsInfo{
		nc: nc,
		js: js,
	}

	return natsInfo, nil
}

func (n *NatsInfo) Publish(msg *nats.Msg) error {

	if msg == nil {
		return errors.New("invalid message")
	}

	printHeaders(msg.Header)

	_, err := n.js.PublishMsg(msg)

	if err != nil {
		return err
	}

	return nil
}

func printHeaders(headers nats.Header) {
	for key, element := range headers {
		zap.S().Debugf("header key [%s], contents [%v]", key, element)
	}
}
