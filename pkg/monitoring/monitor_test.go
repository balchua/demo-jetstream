package monitoring

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/balchua/demo-jetstream/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const zeroLagFakeResponse = `{
    "server_id": "NBKARLLZ7VQLMV2WYSI3TORTBT3XHJUPFXSHMPF4RHW4W4YBXMYQRTPC",
    "now": "2021-12-30T02:09:33.701593895Z",
    "config": {
        "max_memory": 1073741824,
        "max_storage": 1073741824,
        "store_dir": "/data/jetstream"
    },
    "memory": 0,
    "storage": 2920,
    "accounts": 1,
    "api": {
        "total": 0,
        "errors": 0
    },
    "current_api_calls": 0,
    "total_streams": 1,
    "total_consumers": 1,
    "total_messages": 20,
    "total_message_bytes": 2920,
    "account_details": [
        {
            "name": "demo",
            "id": "demo",
            "memory": 0,
            "storage": 2920,
            "api": {
                "total": 6,
                "errors": 0
            },
            "stream_detail": [
                {
                    "name": "USER_TXN",
                    "cluster": {
                        "leader": "bnats-0"
                    },
                    "state": {
                        "messages": 20,
                        "bytes": 2920,
                        "first_seq": 1,
                        "first_ts": "2021-12-30T02:06:05.287010866Z",
                        "last_seq": 20,
                        "last_ts": "2021-12-30T02:07:28.63796658Z",
                        "consumer_count": 1
                    },
                    "consumer_detail": [
                        {
                            "stream_name": "USER_TXN",
                            "name": "GRP_MAKER",
                            "created": "2021-12-30T02:04:27.636237752Z",
                            "delivered": {
                                "consumer_seq": 20,
                                "stream_seq": 20,
                                "last_active": "2021-12-30T02:07:28.637989535Z"
                            },
                            "ack_floor": {
                                "consumer_seq": 20,
                                "stream_seq": 20,
                                "last_active": "2021-12-30T02:07:57.762392703Z"
                            },
                            "num_ack_pending": 0,
                            "num_redelivered": 0,
                            "num_waiting": 1,
                            "num_pending": 1,
                            "cluster": {
                                "leader": "bnats-0"
                            }
                        }
                    ]
                }
            ]
        }
    ]
}`

type MonitoringTestSuite struct {
	suite.Suite
	logs *observer.ObservedLogs
}

func (testSuite *MonitoringTestSuite) SetupTest() {
	var observedZapCore zapcore.Core
	observedZapCore, testSuite.logs = observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	zap.ReplaceGlobals(observedLogger)
}

func (testSuite *MonitoringTestSuite) TestDoStuffWithTestServer() {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(zeroLagFakeResponse))
	}))
	// Close the server when test finishes
	defer server.Close()

	fmt.Printf("Server URL %s", server.URL)
	c := server.Client()

	mConfig := config.Monitor{
		Account:      "demo",
		ConsumerName: "GRP_MAKER",
		StreamName:   "USER_TXN",
		PollSeconds:  3,
	}
	// c := http.Client{Timeout: time.Duration(1) * time.Second}
	m := NewMonitor(mConfig, c, server.URL)
	go m.StartMonitor()
	time.Sleep(3 * time.Second)

	var logExist bool
	logExist = false
	appLogs := testSuite.logs.All()
	for _, appLog := range appLogs {
		fmt.Printf("log content: %s\n", appLog.Message)
		if strings.Contains(appLog.Message, "total lag is 1") {
			logExist = true
		}
	}

	assert.Equal(testSuite.T(), true, logExist)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(MonitoringTestSuite))
}
