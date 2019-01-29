package model

import (
	"github.com/cloudflare/cloudflare-go"
	"time"
)

type Aggregate struct {
	ZoneName string
	ZoneID   string

	Totals map[time.Time]*Counters
}

type Counters struct {
	RequestAll        Counter
	RequestCached     Counter
	RequestUncached   Counter
	BandwidthAll      Counter
	BandwidthCached   Counter
	BandwidthUncached Counter
	HTTPStatus        map[string]Counter
	WafTrigger        map[string]*WafTriggerCounters
}

type WafTriggerCounters struct {
	Simulate    Counter
	Block       Counter
	Challenge   Counter
	JSChallenge Counter
}

type Counter struct {
	Key   string
	Value int
}

func NewCounters() *Counters {
	return &Counters{
		RequestAll:        Counter{Key: "total.requests.all", Value: 0},
		RequestCached:     Counter{Key: "total.requests.cached", Value: 0},
		RequestUncached:   Counter{Key: "total.requests.uncached", Value: 0},
		BandwidthAll:      Counter{Key: "total.bandwidth.all", Value: 0},
		BandwidthCached:   Counter{Key: "total.bandwidth.cached", Value: 0},
		BandwidthUncached: Counter{Key: "total.bandwidth.uncached", Value: 0},
		HTTPStatus: map[string]Counter{
			"2xx": {Key: "total.requests.http_status.2xx", Value: 0},
			"3xx": {Key: "total.requests.http_status.3xx", Value: 0},
			"4xx": {Key: "total.requests.http_status.4xx", Value: 0},
			"5xx": {Key: "total.requests.http_status.5xx", Value: 0}},
		WafTrigger: map[string]*WafTriggerCounters{},
	}
}

func NewWafTriggerResult() *WafTriggerCounters {
	return &WafTriggerCounters{
		Simulate:    Counter{Key: "total.waf.trigger.simulate", Value: 0},
		Block:       Counter{Key: "total.waf.trigger.block", Value: 0},
		Challenge:   Counter{Key: "total.waf.trigger.challenge", Value: 0},
		JSChallenge: Counter{Key: "total.waf.trigger.jschallenge", Value: 0},
	}
}

func NewAggregate(zone cloudflare.Zone) *Aggregate {
	return &Aggregate{
		ZoneName: zone.Name,
		ZoneID:   zone.ID,
		Totals:   make(map[time.Time]*Counters, 0),
	}
}
