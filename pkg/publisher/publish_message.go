package publisher

import (
	"context"

	"github.com/balchua/demo-jetstream/pkg/dtrace"
	"github.com/balchua/demo-jetstream/pkg/infra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UserTransactionPublisher struct {
	natsInfo infra.Nats
}

func NewTransactionPublisher(natsInfo infra.Nats) *UserTransactionPublisher {

	return &UserTransactionPublisher{
		natsInfo: natsInfo,
	}
}

func (u *UserTransactionPublisher) SendMessage(ctx context.Context, message string, subject string) error {
	var span trace.Span
	propagator := otel.GetTextMapPropagator()

	ctx, span = otel.Tracer("publisher").Start(ctx, "SendMessage")
	defer span.End()
	zap.S().Infof("%s", message)
	m := infra.NewNatsMessage(subject)
	carrier := dtrace.NewNatsMessageCarrier(m)
	propagator.Inject(ctx, carrier)
	m.AddHeader("CUSTOM_HEADER", "user-txn")

	m.SetBody([]byte(message))
	err := u.natsInfo.Publish(m)
	if err != nil {
		return err
	}

	return nil

}
