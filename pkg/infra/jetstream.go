package infra

import (
	"github.com/nats-io/jsm.go"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type JetStreamInfo struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewJetStream(seedFile string, natsUri string) (*JetStreamInfo, error) {
	var err error
	var nc *nats.Conn
	var js nats.JetStreamContext
	opt, err := nats.NkeyOptionFromSeed(seedFile)

	if err != nil {
		return nil, err
	}

	nc, err = nats.Connect(natsUri, opt)

	if err != nil {
		return nil, err
	}

	js, err = nc.JetStream()

	if err != nil {
		return nil, err
	}

	return &JetStreamInfo{
		nc: nc,
		js: js,
	}, nil
}

func (jetstream *JetStreamInfo) CreateStream(streamName string, streamOpts ...jsm.StreamOption) error {
	mgr, err := jsm.New(jetstream.nc)

	if err != nil {
		return err
	}
	_, streamErr := mgr.NewStream(streamName, streamOpts...)

	if streamErr != nil && streamErr.Error() != "stream name already in use" {
		return nil
	}
	return nil
}

func (jetstream *JetStreamInfo) CreateConsumer(streamName string, consumerOpts ...jsm.ConsumerOption) error {
	mgr, _ := jsm.New(jetstream.nc)
	_, err := mgr.NewConsumer(streamName, consumerOpts...)

	if err != nil {
		zap.S().Errorf("%v", err)
	}

	return nil
}
