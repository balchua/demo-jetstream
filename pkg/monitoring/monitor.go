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

func Monitor(monitorConfig config.Monitor) {
	for {
		doMonitor(monitorConfig)
	}

}

func doMonitor(mc config.Monitor) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	uri := fmt.Sprintf("%s://%s:%d/jsz?consumers=true", mc.Scheme, mc.Host, mc.Port)
	resp, err := c.Get(uri)
	if err != nil {
		zap.S().Errorf("%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	jsz := &JSInfo{}
	json.Unmarshal(body, jsz)
	checkPendingMessages(jsz)
	time.Sleep(time.Duration(mc.PollSeconds) * time.Second)
}

func checkPendingMessages(jsz *JSInfo) {

	for _, d := range jsz.AccountDetails {
		if d.Name == "demo" {
			if d.Streams != nil {
				for _, s := range d.Streams {
					for _, cons := range s.Consumers {
						if cons.Stream == "USER_TXN" && cons.Name == "GRP_MAKER" {
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
