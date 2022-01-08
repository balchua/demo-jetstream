package publisher

import (
	"context"
	"errors"
	"testing"

	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type PublisherTestSuite struct {
	suite.Suite
	logs         *observer.ObservedLogs
	streamName   string
	consumerName string
}

func (testSuite *PublisherTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, testSuite.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
	testSuite.streamName = "TEST"
	testSuite.consumerName = "TEST_CONSUMER"
}

func (testSuite *PublisherTestSuite) TestPublishMessage() {

	testMessage := `{"TransactionID":1,"UserID":1,"Status":"OK","Amount": 456.89}`
	mockNats := new(infra.MockNats)
	pub := NewTransactionPublisher(mockNats)
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(testMessage))

	mockNats.On("Publish", mock.Anything).Return(nil)

	err := pub.SendMessage(context.Background(), testMessage, testSuite.streamName)
	assert.Nil(testSuite.T(), err)
}

func (testSuite *PublisherTestSuite) TestFailToPublishMessage() {

	testMessage := `{"TransactionID":1,"UserID":1,"Status":"OK","Amount": 456.89}`
	mockNats := new(infra.MockNats)
	pub := NewTransactionPublisher(mockNats)
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(testMessage))

	mockNats.On("Publish", mock.Anything).Return(errors.New("fail to publish"))

	err := pub.SendMessage(context.Background(), testMessage, testSuite.streamName)
	assert.NotNil(testSuite.T(), err)
	assert.Equal(testSuite.T(), "fail to publish", err.Error())
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(PublisherTestSuite))
}
