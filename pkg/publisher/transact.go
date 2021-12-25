package publisher

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type UserTransactionPublisher struct {
	js nats.JetStreamContext
}

func NewTransactionPublisher() *UserTransactionPublisher {
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
