package infra

import (
	"time"

	"github.com/nats-io/jsm.go"
	"go.uber.org/zap"
)

func SetupInfra(natsInfo *NatsInfo, streamSubjects string, streamName string) error {

	consumerName := "GRP_MAKER"
	consumerSubject := "USER_TXN.maker"
	zap.S().Infof("setting up stream [%s] stream subject [%s]", streamName, streamSubjects)
	if err := setupStream(natsInfo, streamSubjects, streamName); err != nil {
		return err
	}

	zap.S().Infof("setting up consumer name [%s] on subject [%s]", consumerName, consumerSubject)
	if err := setupConsumer(natsInfo, streamName, consumerName, consumerSubject); err != nil {
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
