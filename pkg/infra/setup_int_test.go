package infra

import (
	"testing"
	"time"

	_ "github.com/balchua/demo-jetstream/pkg/test_util"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type SetupInfraTestSuite struct {
	suite.Suite
	logs            *observer.ObservedLogs
	streamName      string
	consumerName    string
	streamSubjects  string
	streamMaxAge    time.Duration
	consumerSubject string
}

// This is an integration test, as we need a real NATS Jetstream running
func (s *SetupInfraTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, s.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
	s.streamName = "USER_TXN"
	s.consumerName = "GRP_MAKER"
	s.streamSubjects = "USER_TXN.>"
	s.streamMaxAge = 1 * time.Minute
	s.consumerSubject = "USER_TXN.maker"
}

func (s *SetupInfraTestSuite) TestITSetupInfra() {
	seedFile := "hack/sys-seed.txt"
	natsUri := "localhost:4220"
	jsi, err := NewJetStream(seedFile, natsUri)

	if err != nil {
		s.Fail("setup failure \n%v", err)
	}

	i := NewInfraSetup(jsi, s.streamName, s.streamSubjects, s.consumerName, s.consumerSubject, 1*time.Minute, 1)
	if err := i.Setup(); err != nil {
		s.Fail("setup failure \n%v", err)
	}
	s.Assert().True(jsi.IsStreamAvailable(s.streamName))
	jsi.Close()
}

func TestSetupInfraSuite(t *testing.T) {
	suite.Run(t, new(SetupInfraTestSuite))
}
