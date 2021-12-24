package infra

import (
	"time"

	"github.com/nats-io/jsm.go"
	"go.uber.org/zap"
)

func SetupInfra(natsInfo *NatsInfo, streamSubjects string, streamName string) error {
	if err := setupStream(natsInfo, streamSubjects, streamName); err != nil {
		return err
	}

	if err := setupConsumer(natsInfo, streamName, "GRP_MAKER", "USER_TXN.maker"); err != nil {
		return err
	}

	return nil
}

func setupStream(natsInfo *NatsInfo, streamSubjects string, streamName string) error {
	mgr, _ := jsm.New(natsInfo.Nc)
	_, err := mgr.NewStream(streamName, jsm.Subjects(streamSubjects), jsm.MaxAge(24*365*time.Hour), jsm.FileStorage())
	// Check if the stream already exists; if not, create it.
	if err != nil {
		zap.S().Errorf("%v", err)
	}
	return nil

}

func setupConsumer(natsInfo *NatsInfo, streamName string, consumerName string, filterSubject string) error {
	mgr, _ := jsm.New(natsInfo.Nc)
	_, err := mgr.NewConsumer(streamName, jsm.DurableName(consumerName),
		jsm.FilterStreamBySubject(filterSubject),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
		jsm.ReplayAsReceived())

	if err != nil {
		zap.S().Errorf("%v", err)
	}

	return nil
}
