package publisher

import (
	"github.com/balchua/demo-jetstream/pkg/infra"
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

func (u *UserTransactionPublisher) SendMessage(message string, subject string) error {
	zap.S().Infof("%s", message)
	m := infra.NewNatsMessage(subject)
	m.AddHeader("CUSTOM_HEADER", "user-txn")

	m.SetBody([]byte(message))
	err := u.natsInfo.Publish(m)
	if err != nil {
		return err
	}
	return nil

}
