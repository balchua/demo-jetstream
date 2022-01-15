package dtrace

import (
	"context"
	"testing"

	"github.com/balchua/demo-jetstream/pkg/config"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest/observer"
)

type DTraceTestSuite struct {
	suite.Suite
	logs *observer.ObservedLogs
}

func (testSuite *NatsMessageCarrierTestSuite) TestDTraceSetup() {

	appConfig := &config.AppConfiguration{}
	appConfig.T = config.Tracing{
		JaegerUrl:   "http://dummy/api/traces",
		ServiceName: "integrationTest",
	}

	dt := SetupTracer(appConfig.T)
	dt.Flush(context.TODO())
	dt.Close()
}

func TestDTraceSuite(t *testing.T) {
	suite.Run(t, new(DTraceTestSuite))
}
