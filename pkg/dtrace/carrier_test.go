package dtrace

import (
	"testing"

	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type NatsMessageCarrierTestSuite struct {
	suite.Suite
	logs    *observer.ObservedLogs
	carrier NatsMessageCarrier
}

func (testSuite *NatsMessageCarrierTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, testSuite.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
	msg := infra.NewNatsMessage("TEST")

	testSuite.carrier = NewNatsMessageCarrier(msg)
	testSuite.carrier.Set("header1", "value1")
	testSuite.carrier.Set("header2", "value2")
}

func (testSuite *NatsMessageCarrierTestSuite) TestCarrierWithExistingKey() {

	value := testSuite.carrier.Get("header1")
	assert.Equal(testSuite.T(), "value1", value)
}

func (testSuite *NatsMessageCarrierTestSuite) TestCarrierWithNoneExistingKey() {

	value := testSuite.carrier.Get("headerX")
	assert.Equal(testSuite.T(), "", value)
}

func (testSuite *NatsMessageCarrierTestSuite) TestCarrierAllKeysExist() {
	values := testSuite.carrier.Keys()

	for value := range values {
		zap.S().Debugf("value is %s", value)
	}
	// assert.Equal(testSuite.T(), "header2", values[0])
	// assert.Equal(testSuite.T(), "header1", values[1])
}
func TestNatsMessageCarrierSuite(t *testing.T) {
	suite.Run(t, new(NatsMessageCarrierTestSuite))
}
