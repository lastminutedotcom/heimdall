package model

import (
	"fmt"
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
	WafTrigger        map[string]*WafActionCounters
	RateLimit         map[string]map[string]*RateLimitCounters
}

type WafActionCounters struct {
	Simulate    Counter
	Block       Counter
	Challenge   Counter
	JSChallenge Counter
}

type RateLimitCounters struct {
	Simulate        Counter
	Drop            Counter
	Challenge       Counter
	JSChallenge     Counter
	ConnectionClose Counter
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
		WafTrigger: map[string]*WafActionCounters{},
		RateLimit:  map[string]map[string]*RateLimitCounters{},
	}
}

func NewWafTriggerResult() *WafActionCounters {
	return &WafActionCounters{
		Simulate:    Counter{Key: "total.waf.trigger.simulate", Value: 0},
		Block:       Counter{Key: "total.waf.trigger.block", Value: 0},
		Challenge:   Counter{Key: "total.waf.trigger.challenge", Value: 0},
		JSChallenge: Counter{Key: "total.waf.trigger.jschallenge", Value: 0},
	}
}

func NewRateLimitResult() map[string]*RateLimitCounters {
	securityEventCounters := make(map[string]*RateLimitCounters, 0)
	securityEventCounters["GET"] = NewSecurityEventCounters("get")
	securityEventCounters["POST"] = NewSecurityEventCounters("post")
	securityEventCounters["PUT"] = NewSecurityEventCounters("put")
	securityEventCounters["PATCH"] = NewSecurityEventCounters("patch")
	securityEventCounters["DELETE"] = NewSecurityEventCounters("delete")
	return securityEventCounters
}

func NewSecurityEventCounters(method string) *RateLimitCounters {
	return &RateLimitCounters{
		Simulate:        Counter{Key: fmt.Sprintf("total.ratelimit.%s.simulate", method), Value: 0},
		Drop:            Counter{Key: fmt.Sprintf("total.ratelimit.%s.drop", method), Value: 0},
		Challenge:       Counter{Key: fmt.Sprintf("total.ratelimit.%s.challenge", method), Value: 0},
		JSChallenge:     Counter{Key: fmt.Sprintf("total.ratelimit.%s.jschallenge", method), Value: 0},
		ConnectionClose: Counter{Key: fmt.Sprintf("total.ratelimit.%s.connection_close", method), Value: 0},
	}
}
func NewAggregate(zone cloudflare.Zone) *Aggregate {
	return &Aggregate{
		ZoneName: zone.Name,
		ZoneID:   zone.ID,
		Totals:   make(map[time.Time]*Counters, 0),
	}
}
