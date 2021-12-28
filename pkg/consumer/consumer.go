package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/balchua/demo-jetstream/pkg/model"
	"github.com/nats-io/nats.go"

	"go.uber.org/zap"
)

type Consumer struct {
	natsInfo infra.Nats
}

func NewConsumer(natsInfo infra.Nats) *Consumer {

	return &Consumer{
		natsInfo: natsInfo,
	}
}

func (c *Consumer) Listen(ctx context.Context, done chan bool, subject string, consumerName string, sleepTimeInMillis int) {
	err := c.natsInfo.Subscribe(subject, consumerName)
	if err != nil {
		zap.S().Fatalf("unable to subscribe to subject %s, %v", subject, err)
	}

	for {
		select {
		case <-ctx.Done():
			zap.S().Info("ready to end worker process")
			done <- true
			return
		default:
		}
		msgs, err := c.natsInfo.Fetch(100, nats.Context(ctx))

		if err != nil {
			zap.S().Fatalf("unable to consume message %v", err)
		}
		for _, msg := range msgs {
			msg.Ack()
			var userTxn model.UserTransaction
			zap.S().Infof("header: [%s]", msg.GetHeader("CUSTOM_HEADER"))
			err := json.Unmarshal(msg.GetBody(), &userTxn)
			zap.S().Infof("%s", msg.GetBody())
			if err != nil {
				zap.S().Errorf("%v", err)
			}
			zap.S().Infof("TransactionId: %d, Amount: %s, Status: %s", userTxn.TransactionID, userTxn.Amount.String(), userTxn.Status)
			time.Sleep(time.Duration(sleepTimeInMillis) * time.Millisecond)
		}
		zap.S().Info("done processing messages from this batch")

	}
}
