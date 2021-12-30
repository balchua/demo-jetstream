package monitoring

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/balchua/demo-jetstream/pkg/config"
	"go.uber.org/zap"
)

type Monitor struct {
	client  *http.Client
	baseURL string
	config  config.Monitor
}

func NewMonitor(monitorConfig config.Monitor, client *http.Client, baseURL string) *Monitor {

	return &Monitor{
		client:  client,
		baseURL: baseURL,
		config:  monitorConfig,
	}
}

func (m *Monitor) StartMonitor() {
	for {
		m.doMonitor(m.config)
	}
}

func (m *Monitor) doMonitor(mc config.Monitor) {
	uri := fmt.Sprintf("%s/jsz?consumers=true", m.baseURL)
	resp, err := m.client.Get(uri)
	if err != nil {
		zap.S().Errorf("%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	jsz := &JSInfo{}
	json.Unmarshal(body, jsz)
	checkPendingMessages(jsz, mc)
	time.Sleep(time.Duration(mc.PollSeconds) * time.Second)
}

func checkPendingMessages(jsz *JSInfo, mc config.Monitor) {

	for _, d := range jsz.AccountDetails {
		if d.Name == mc.Account {
			if d.Streams != nil {
				for _, s := range d.Streams {
					for _, cons := range s.Consumers {
						if cons.Stream == mc.StreamName && cons.Name == mc.ConsumerName {
							totalLag := cons.NumPending + cons.NumAckPending
							if totalLag > 0 {
								zap.S().Infof("total lag is %d", totalLag)
								break
							}
						}
					}
				}
			}
		}

	}
}
