package model

import (
	"github.com/cloudflare/cloudflare-go"
	"time"
)

type Aggregate struct {
	ZoneName string
	ZoneID   string

	Date                   time.Time
	TotalRequestAll        KeyValue
	TotalRequestCached     KeyValue
	TotalRequestUncached   KeyValue
	TotalBandwidthAll      KeyValue
	TotalBandwidthCached   KeyValue
	TotalBandwidthUncached KeyValue
	HTTPStatus             map[string]KeyValue
}

type KeyValue struct {
	Key   string
	Value int
}

func NewAggregate(zone cloudflare.Zone) *Aggregate {
	return &Aggregate{
		ZoneName:               zone.Name,
		ZoneID:                 zone.ID,
		TotalRequestAll:        KeyValue{Key: "total.requests.all", Value: 0},
		TotalRequestCached:     KeyValue{Key: "total.requests.cached", Value: 0},
		TotalRequestUncached:   KeyValue{Key: "total.requests.uncached", Value: 0},
		TotalBandwidthAll:      KeyValue{Key: "total.bandwidth.all", Value: 0},
		TotalBandwidthCached:   KeyValue{Key: "total.bandwidth.cached", Value: 0},
		TotalBandwidthUncached: KeyValue{Key: "total.bandwidth.uncached", Value: 0},
		HTTPStatus: map[string]KeyValue{
			"2xx": {Key: "total.requests.http_status.2xx", Value: 0},
			"3xx": {Key: "total.requests.http_status.3xx", Value: 0},
			"4xx": {Key: "total.requests.http_status.4xx", Value: 0},
			"5xx": {Key: "total.requests.http_status.5xx", Value: 0}},
	}
}
