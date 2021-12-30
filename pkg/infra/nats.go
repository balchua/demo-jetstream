package infra

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NatsImpl struct {
	nc  *nats.Conn
	js  nats.JetStreamContext
	sub *nats.Subscription
}

func NewNats(seedFile string, natsUri string) (*NatsImpl, error) {
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

	natsInfo := &NatsImpl{
		nc: nc,
		js: js,
	}

	return natsInfo, nil
}

func (n *NatsImpl) Publish(msg *NatsMessage) error {

	if msg == nil {
		return errors.New("invalid message")
	}

	printHeaders(msg.GetHeaders())

	_, err := n.js.PublishMsg(msg.GetUnderlyingNatsMessage())

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

func (n *NatsImpl) Subscribe(subject, consumerName string) error {
	var err error
	n.sub, err = n.js.PullSubscribe(subject, consumerName)
	if err != nil {
		return err
	}
	return nil
}

func (n *NatsImpl) Fetch(messageCount int, ctx context.Context) ([]*NatsMessage, error) {
	var natsMessages []*NatsMessage
	msgs, err := n.sub.Fetch(messageCount, nats.Context(ctx))

	if err != nil {
		if !errors.Is(err, nats.ErrContextAndTimeout) &&
			!errors.Is(err, context.DeadlineExceeded) &&
			!errors.Is(err, nats.ErrBadSubscription) &&
			ctx.Err() != context.Canceled {
			zap.S().Errorf("%v", err)
			return nil, err
		}

	}
	for _, msg := range msgs {
		natsMsg := &NatsMessage{
			msg: msg,
		}
		natsMessages = append(natsMessages, natsMsg)
	}
	return natsMessages, nil

}

func (n *NatsImpl) Close() {
	n.nc.Close()
}
