package publisher

import (
	"encoding/json"

	"github.com/balchua/demo-jetstream/pkg/model"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type UserTransactionPublisher struct {
	js nats.JetStreamContext
}

func NewTransactionPublisher() *UserTransactionPublisher {
	nc, _ := nats.Connect("localhost:32422")
	js, err := nc.JetStream()
	if err != nil {
		zap.S().Errorf("%v", err)
	}

	return &UserTransactionPublisher{
		js: js,
	}
}

func (u *UserTransactionPublisher) Publish(message string, subject string) error {

	usr := model.UserTransaction{
		TransactionID: 1,
		UserID:        1,
		Status:        "OK",
	}
	usrJson, _ := json.Marshal(usr)
	_, err := u.js.Publish(subject, usrJson)
	if err != nil {
		return err
	}
	return nil

}
