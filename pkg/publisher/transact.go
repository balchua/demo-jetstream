package publisher

import (
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

	// usr := model.UserTransaction{
	// 	TransactionID: 1,
	// 	UserID:        1,
	// 	Status:        "OK",
	// }
	// usrJson, _ := json.Marshal(usr)
	zap.S().Infof("%s", message)
	m := nats.NewMsg(subject)

	m.Header.Add("CUSTOM_HEADER", "user-txn")
	m.Data = []byte(message)
	_, err := u.js.PublishMsg(m)
	if err != nil {
		return err
	}
	return nil

}
