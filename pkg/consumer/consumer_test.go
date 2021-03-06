package consumer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/balchua/demo-jetstream/pkg/api/verifier"
	"github.com/balchua/demo-jetstream/pkg/infra"
	"github.com/balchua/demo-jetstream/pkg/metrics"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type ConsumerTestSuite struct {
	suite.Suite
	logs         *observer.ObservedLogs
	streamName   string
	consumerName string
	request      *verifier.VerifyTransactionRequest
	response     *verifier.VerifyTransactionResponse
	appMetrics   *metrics.Metrics
}

func (testSuite *ConsumerTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, testSuite.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
	testSuite.appMetrics = metrics.NewMetrics()
	testSuite.streamName = "TEST"
	testSuite.consumerName = "TEST_CONSUMER"
	tx := &verifier.Transaction{
		Status:        "OK",
		TransactionID: 1,
		Amount:        "459.89",
		UserId:        1,
	}
	testSuite.request = &verifier.VerifyTransactionRequest{
		Tx: tx,
	}

	testSuite.response = &verifier.VerifyTransactionResponse{
		Tx:      tx,
		Code:    verifier.StatusCode_OK,
		Message: "mock response",
	}
}
func (testSuite *ConsumerTestSuite) TestFetchMessage() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	mockNats := new(infra.MockNats)
	mockApi := new(verifier.MockTransactionVerifier)
	mockApi.On("VerifyTransaction", mock.Anything, mock.Anything).Return(testSuite.response, nil)
	mockApi.On("Close").Return(testSuite.response)

	con := NewConsumer(mockNats, testSuite.appMetrics, mockApi)

	var messages []*infra.NatsMessage
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(`{"TransactionID":1,"UserID":1,"Status":"OK","Amount": 456.89}`))

	messages = append(messages, msg)

	mockNats.On("Subscribe", testSuite.streamName, testSuite.consumerName).Return(nil)
	mockNats.On("Fetch", 10, ctx).Return(messages, nil)

	worker := make(chan bool)
	go con.Listen(ctx, worker, testSuite.streamName, testSuite.consumerName, 100)
	time.Sleep(10 * time.Millisecond)

	cancelFunc()
	<-worker

	var logExist bool
	logExist = false
	appLogs := testSuite.logs.All()
	for _, appLog := range appLogs {
		fmt.Printf("log content: %s\n", appLog.Message)
		if strings.Contains(appLog.Message, "TransactionId: 1, Amount: 456.89, Status: OK") {
			logExist = true
		}
	}

	mockNats.AssertNumberOfCalls(testSuite.T(), "Fetch", 1)

	assert.Equal(testSuite.T(), true, logExist)
}

func (testSuite *ConsumerTestSuite) TestFailToSubscribe() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	mockNats := new(infra.MockNats)
	mockApi := new(verifier.MockTransactionVerifier)
	mockApi.On("VerifyTransaction", mock.Anything, mock.Anything).Return(testSuite.response, nil)
	mockApi.On("Close").Return(testSuite.response)

	con := NewConsumer(mockNats, testSuite.appMetrics, mockApi)

	var messages []*infra.NatsMessage
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(`{"TransactionID":1,"UserID":1,"Status":"OK","Amount": 456.89}`))

	messages = append(messages, msg)

	mockNats.On("Subscribe", testSuite.streamName, testSuite.consumerName).Return(errors.New("invalid subscription"))

	worker := make(chan bool)
	go con.Listen(ctx, worker, testSuite.streamName, testSuite.consumerName, 100)
	time.Sleep(10 * time.Millisecond)

	cancelFunc()
	<-worker

	var logExist bool
	logExist = false
	appLogs := testSuite.logs.All()
	for _, appLog := range appLogs {
		fmt.Printf("log content: %s\n", appLog.Message)
		if strings.Contains(appLog.Message, "invalid subscription") {
			logExist = true
		}
	}
	mockNats.AssertNumberOfCalls(testSuite.T(), "Subscribe", 1)
	mockNats.AssertNumberOfCalls(testSuite.T(), "Fetch", 0)
	assert.Equal(testSuite.T(), true, logExist)
}

func (testSuite *ConsumerTestSuite) TestFailToFetch() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	mockNats := new(infra.MockNats)
	mockApi := new(verifier.MockTransactionVerifier)
	mockApi.On("VerifyTransaction", mock.Anything, mock.Anything).Return(testSuite.response, nil)
	mockApi.On("Close").Return(testSuite.response)

	con := NewConsumer(mockNats, testSuite.appMetrics, mockApi)

	var messages []*infra.NatsMessage
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(`{"TransactionID":1,"UserID":1,"Status":"OK","Amount": 456.89}`))

	messages = append(messages, msg)

	mockNats.On("Subscribe", testSuite.streamName, testSuite.consumerName).Return(nil)
	mockNats.On("Fetch", 10, ctx).Return(messages, errors.New("fetch data failed"))

	worker := make(chan bool)
	go con.Listen(ctx, worker, testSuite.streamName, testSuite.consumerName, 100)
	time.Sleep(10 * time.Millisecond)

	cancelFunc()
	<-worker

	var logExist bool
	logExist = false
	appLogs := testSuite.logs.All()
	for _, appLog := range appLogs {
		fmt.Printf("log content: %s\n", appLog.Message)
		if strings.Contains(appLog.Message, "fetch data failed") {
			logExist = true
		}
	}
	mockNats.AssertNumberOfCalls(testSuite.T(), "Fetch", 1)
	assert.Equal(testSuite.T(), true, logExist)
}

func (testSuite *ConsumerTestSuite) TestFetchInvalidMessage() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	mockNats := new(infra.MockNats)
	mockApi := new(verifier.MockTransactionVerifier)
	mockApi.On("VerifyTransaction", mock.Anything, mock.Anything).Return(testSuite.response, nil)
	mockApi.On("Close").Return(testSuite.response)

	con := NewConsumer(mockNats, testSuite.appMetrics, mockApi)

	var messages []*infra.NatsMessage
	msg := infra.NewNatsMessage("TEST")

	msg.SetBody([]byte(`{"TransactionID:1,"UserID":1,"Status":"OK","Amount": 456.89}`))

	messages = append(messages, msg)

	mockNats.On("Subscribe", testSuite.streamName, testSuite.consumerName).Return(nil)
	mockNats.On("Fetch", 10, ctx).Return(messages, nil)

	worker := make(chan bool)
	go con.Listen(ctx, worker, testSuite.streamName, testSuite.consumerName, 100)
	time.Sleep(10 * time.Millisecond)

	cancelFunc()
	<-worker

	var logExist bool
	logExist = false
	appLogs := testSuite.logs.All()
	for _, appLog := range appLogs {
		fmt.Printf("log content: %s\n", appLog.Message)
		if strings.Contains(appLog.Message, "invalid character") {
			logExist = true
		}
	}
	mockNats.AssertNumberOfCalls(testSuite.T(), "Fetch", 1)
	mockNats.AssertNumberOfCalls(testSuite.T(), "Subscribe", 1)
	assert.Equal(testSuite.T(), true, logExist)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}
