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

type NatsTestSuite struct {
	suite.Suite
	logs            *observer.ObservedLogs
	streamName      string
	consumerName    string
	streamSubjects  string
	streamMaxAge    time.Duration
	consumerSubject string
}

func (n *NatsTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, n.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
	n.streamName = "USER_TXN"
	n.consumerName = "GRP_MAKER"
	n.streamSubjects = "USER_TXN.>"
	n.streamMaxAge = 1 * time.Minute
	n.consumerSubject = "USER_TXN.maker"
	n.createStream()
}

func (n *NatsTestSuite) TestPublishMessage() {

	seedFile := "hack/seed.txt"
	natsUri := "localhost:32422"
	natsInfo, err := NewNats(seedFile, natsUri)
	if err != nil {
		n.Fail("unable to connect to nats server\n %v", err)
	}
	msgString := `{"TransactionID":1,"UserID":1,"Status":"OK","Amount": 456.89}`
	msg := NewNatsMessage("USER_TXN.maker")
	msg.AddHeader("test", "halo")
	msg.SetBody([]byte(msgString))
	if err := natsInfo.Publish(msg); err != nil {
		n.Fail("unable to publish message %v", err)
	}
	natsInfo.Close()
}
func (s *NatsTestSuite) createStream() {

	seedFile := "hack/sys-seed.txt"
	natsUri := "localhost:32422"
	jsi, err := NewJetStream(seedFile, natsUri)
	if err != nil {
		s.Fail("unable to connect to jetstream \n%v", err)
	}
	if err := jsi.CreateStream(s.streamName, jsm.Subjects(s.streamSubjects), jsm.MaxAge(s.streamMaxAge), jsm.FileStorage()); err != nil {
		s.Fail("unable to create stream in jetstream \n%v", err)
	}
	isAvail, err := jsi.IsStreamAvailable(s.streamName)
	s.Assert().True(isAvail)
	jsi.Close()
}

func (s *JetstreamTestSuite) createConsumer() {

	seedFile := "hack/sys-seed.txt"
	natsUri := "localhost:32422"
	jsi, err := NewJetStream(seedFile, natsUri)
	if err != nil {
		s.Fail("unable to connect to jetstream \n%v", err)
	}
	if err := jsi.CreateConsumer(s.streamName, jsm.DurableName(s.consumerName),
		jsm.FilterStreamBySubject(s.consumerSubject),
		jsm.AcknowledgeExplicit(),
		jsm.DeliverAllAvailable(),
		jsm.ReplayAsReceived(),
		jsm.MaxWaiting(1)); err != nil {
		s.Fail("unable to connect to jetstream \n%v", err)
	}
	jsi.Close()
}

func TestPublishIntegration(t *testing.T) {
	suite.Run(t, new(NatsTestSuite))
}
