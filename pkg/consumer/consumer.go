package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/balchua/demo-jetstream/pkg/dtrace"
	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/balchua/demo-jetstream/pkg/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

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
		zap.S().Errorf("unable to subscribe to subject %s, %v", subject, err)
		done <- true
		return
	}

	for {
		select {
		case <-ctx.Done():
			zap.S().Info("ready to end worker process")
			done <- true
			return
		default:
		}
		msgs, err := c.natsInfo.Fetch(1, ctx)

		if err != nil {
			zap.S().Errorf("unable to consume message %v", err)
			done <- true
			return
		}
		for _, msg := range msgs {
			msg.Ack()
			var userTxn model.UserTransaction
			propagator := otel.GetTextMapPropagator()
			carrier := dtrace.NewNatsMessageCarrier(msg)
			ctx = propagator.Extract(ctx, carrier)
			var span trace.Span
			ctx, span = otel.Tracer("listener").Start(ctx, "Listen")
			zap.S().Infof("header: [%s]", msg.GetHeader("CUSTOM_HEADER"))
			zap.S().Infof("header: [%s]", msg.GetHeader("traceparent"))
			err := json.Unmarshal(msg.GetBody(), &userTxn)
			zap.S().Infof("%s", msg.GetBody())
			if err != nil {
				zap.S().Errorf("%v", err)
			}
			zap.S().Infof("TransactionId: %d, Amount: %s, Status: %s", userTxn.TransactionID, userTxn.Amount.String(), userTxn.Status)
			time.Sleep(time.Duration(sleepTimeInMillis) * time.Millisecond)
			span.End()
		}
		zap.S().Debug("done processing messages from this batch")

	}
}
