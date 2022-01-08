package dtrace

import (
	"github.com/balchua/demo-jetstream/pkg/infra"
)

type NatsMessageCarrier struct {
	msg *infra.NatsMessage
}

func NewNatsMessageCarrier(msg *infra.NatsMessage) NatsMessageCarrier {
	return NatsMessageCarrier{msg: msg}
}

func (p NatsMessageCarrier) Get(key string) string {
	for k, value := range p.msg.GetHeaders() {
		if string(k) == key {
			return string(value[0])
		}
	}
	return ""

}

func (p NatsMessageCarrier) Set(key, val string) {
	p.msg.AddHeader(key, val)
}

func (p NatsMessageCarrier) Keys() []string {
	out := make([]string, len(p.msg.GetHeaders()))
	var i int = 0
	for k, _ := range p.msg.GetHeaders() {
		out[i] = string(k)
	}
	return out
}
