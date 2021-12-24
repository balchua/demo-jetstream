package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/balchua/demo-jetstream/pkg/model"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Consumer struct {
	js nats.JetStreamContext
}

func NewConsumer(natsInfo *infra.NatsInfo) *Consumer {
	nc, _ := nats.Connect("localhost:32422")
	js, err := nc.JetStream()
	if err != nil {
		zap.S().Errorf("%v", err)
	}

	return &Consumer{
		js: js,
	}
}

func (c *Consumer) Listen(ctx context.Context, done chan bool, subject string, consumerName string) {
	sub, _ := c.js.PullSubscribe(subject, consumerName)

	for {
		select {
		case <-ctx.Done():
			zap.S().Info("ready to end worker process")
			done <- true
			return
		default:
		}
		msgs, err := sub.Fetch(1, nats.Context(ctx))
		if err != nil {
			if !errors.Is(err, nats.ErrContextAndTimeout) && !errors.Is(err, context.DeadlineExceeded) && ctx.Err() != context.Canceled {
				zap.S().Errorf("%v", err)
			}

		} else {
			for _, msg := range msgs {
				msg.Ack()
				var userTxn model.UserTransaction
				err := json.Unmarshal(msg.Data, &userTxn)
				zap.S().Infof("%s", msg.Data)
				if err != nil {
					zap.S().Errorf("%v", err)
				}
				zap.S().Infof("TransactionId: %d, Amount: %s, Status: %s", userTxn.TransactionID, userTxn.Amount.String(), userTxn.Status)
				time.Sleep(3 * time.Second)
			}
			zap.S().Info("done processing messages from this batch")
		}

	}
}
