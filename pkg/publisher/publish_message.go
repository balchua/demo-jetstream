package publisher

import (
	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/nats-io/nats.go"
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
	m := nats.NewMsg(subject)

	m.Header.Add("CUSTOM_HEADER", "user-txn")
	m.Data = []byte(message)
	err := u.natsInfo.Publish(m)
	if err != nil {
		return err
	}
	return nil

}
