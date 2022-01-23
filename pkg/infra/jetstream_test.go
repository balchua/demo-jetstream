package infra

import (
	"testing"
	"time"

	_ "github.com/balchua/demo-jetstream/pkg/test_util"
	"github.com/nats-io/jsm.go"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type JetstreamTestSuite struct {
	suite.Suite
	logs            *observer.ObservedLogs
	streamName      string
	consumerName    string
	streamSubjects  string
	streamMaxAge    time.Duration
	consumerSubject string
}

func (j *JetstreamTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, j.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
	j.streamName = "USER_TXN"
	j.consumerName = "GRP_MAKER"
	j.streamSubjects = "USER_TXN.>"
	j.streamMaxAge = 1 * time.Minute
	j.consumerSubject = "USER_TXN.maker"
}

func (j *JetstreamTestSuite) TestCreateStream() {

	seedFile := "hack/sys-seed.txt"
	natsUri := "localhost:4220"
	jsi, err := NewJetStream(seedFile, natsUri)
	if err != nil {
		j.Fail("unable to connect to jetstream \n%v", err)
	}
	if err := jsi.CreateStream(j.streamName, jsm.Subjects(j.streamSubjects), jsm.MaxAge(j.streamMaxAge), jsm.FileStorage()); err != nil {
		j.Fail("unable to create stream in jetstream \n%v", err)
	}
	isAvail, err := jsi.IsStreamAvailable(j.streamName)
	j.Assert().True(isAvail)
	jsi.Close()
}

func (j *JetstreamTestSuite) TestCreateConsumer() {

	seedFile := "hack/sys-seed.txt"
	natsUri := "localhost:4220"
	jsi, err := NewJetStream(seedFile, natsUri)
	if err != nil {
		j.Fail("unable to connect to jetstream \n%v", err)
	}
	if err := jsi.CreateConsumer(j.streamName, jsm.DurableName(j.consumerName),
		jsm.FilterStreamBySubject(j.consumerSubject),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
		jsm.ReplayAsReceived(),
		jsm.MaxWaiting(1)); err != nil {
		j.Fail("unable to connect to jetstream \n%v", err)
	}
	jsi.Close()
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(JetstreamTestSuite))
}
