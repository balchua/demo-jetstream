package infra

import (
	"time"

	"github.com/nats-io/jsm.go"
	"go.uber.org/zap"
)

type InfraSetup struct {
	jetstream       JetStream
	streamName      string
	streamSubjects  string
	consumerName    string
	consumerSubject string
	streamMaxAge    time.Duration
}

func NewInfraSetup(jetstream JetStream, streamName, streamSubjects, consumerName, consumerSubject string, maxAge time.Duration) *InfraSetup {

	zap.S().Infof("setting up stream [%s] stream subject [%s]", streamName, streamSubjects)
	zap.S().Infof("setting up consumer name [%s] on subject [%s]", consumerName, consumerSubject)
	return &InfraSetup{
		jetstream:       jetstream,
		streamName:      streamName,
		streamSubjects:  streamSubjects,
		consumerName:    consumerName,
		consumerSubject: consumerSubject,
		streamMaxAge:    maxAge,
	}
}
func (i *InfraSetup) Setup() error {

	if err := i.jetstream.CreateStream(i.streamName, jsm.Subjects(i.streamSubjects), jsm.MaxAge(i.streamMaxAge), jsm.FileStorage()); err != nil {
		return err
	}

	if err := i.jetstream.CreateConsumer(i.streamName, jsm.DurableName(i.consumerName),
		jsm.FilterStreamBySubject(i.consumerSubject),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
		jsm.ReplayAsReceived(),
		jsm.MaxWaiting(1)); err != nil {
		return err
	}

	return nil
}
