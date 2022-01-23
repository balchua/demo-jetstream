package consumer

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/balchua/demo-jetstream/pkg/api/verifier"
	"github.com/balchua/demo-jetstream/pkg/dtrace"
	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/balchua/demo-jetstream/pkg/metrics"
	"github.com/balchua/demo-jetstream/pkg/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
)

type Consumer struct {
	natsInfo       infra.Nats
	verifierClient verifier.ITransactionVerifier
	appMetrics     *metrics.Metrics
}

func NewConsumer(natsInfo infra.Nats, appMetrics *metrics.Metrics, verifierClient verifier.ITransactionVerifier) *Consumer {

	return &Consumer{
		natsInfo:       natsInfo,
		verifierClient: verifierClient,
		appMetrics:     appMetrics,
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
			c.verifierClient.Close()
			done <- true
			return
		default:
		}
		msgs, err := c.natsInfo.Fetch(10, ctx)

		if err != nil {
			zap.S().Errorf("unable to consume message %v", err)
			done <- true
			return
		}
		for _, msg := range msgs {
			start := time.Now()
			rand.Seed(time.Now().UnixNano())
			msg.Ack()
			var userTxn model.UserTransaction
			propagator := otel.GetTextMapPropagator()
			carrier := dtrace.NewNatsMessageCarrier(msg)
			ctx = propagator.Extract(ctx, carrier)
			var span trace.Span
			var spanctx context.Context
			spanctx, span = otel.Tracer("listener").Start(ctx, "Listen")
			zap.S().Infof("header: [%s]", msg.GetHeader("CUSTOM_HEADER"))
			zap.S().Infof("header: [%s]", msg.GetHeader("traceparent"))
			err := json.Unmarshal(msg.GetBody(), &userTxn)
			zap.S().Infof("%s", msg.GetBody())
			if err != nil {
				zap.S().Errorf("%v", err)
			}
			zap.S().Infof("TransactionId: %d, Amount: %s, Status: %s", userTxn.TransactionID, userTxn.Amount.String(), userTxn.Status)
			tx := &verifier.Transaction{
				UserId:        int64(userTxn.UserID),
				TransactionID: int64(userTxn.TransactionID),
				Amount:        userTxn.Amount.String(),
				Status:        userTxn.Status,
			}

			response, responseErr := c.verifierClient.VerifyTransaction(spanctx, &verifier.VerifyTransactionRequest{
				Tx: tx,
			})

			time.Sleep(time.Duration(rand.Intn(sleepTimeInMillis)) * time.Millisecond)
			if responseErr == nil {
				zap.S().Infof("Response status code [%d] with message [%s]", response.Code, response.Message)
			}

			span.End()
			processDuration := time.Since(start)
			zap.S().Debugf("elapsed: %f", processDuration.Seconds())
			c.appMetrics.Observe(processDuration)
		}
		zap.S().Debug("done processing messages from this batch")

	}
}
