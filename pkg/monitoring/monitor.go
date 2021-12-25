package monitoring

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Monitor(host string, port int, scheme string) {
	for {
		doMonitor(host, port, scheme)
	}

}

func doMonitor(host string, port int, scheme string) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	uri := fmt.Sprintf("%s://%s:%d/jsz?consumers=true", scheme, host, port)
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
		time.Sleep(1000 * time.Millisecond)
	}
}
