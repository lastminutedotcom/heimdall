package model

import (
	"github.com/cloudflare/cloudflare-go"
	"time"
)

type Aggregate struct {
	ZoneName string
	ZoneID   string
	Date     time.Time

	TotalRequestAll        KeyValue
	TotalRequestCached     KeyValue
	TotalRequestUncached   KeyValue
	TotalBandwidthAll      KeyValue
	TotalBandwidthCached   KeyValue
	TotalBandwidthUncached KeyValue

	HTTPStatus map[string]KeyValue
}

type KeyValue struct {
	Key   string
	Value int
}

type ZoneAnalyticsColocationResponse struct {
	cloudflare.Response
	Query struct {
		Since     time.Time `json:"since"`
		Until     time.Time `json:"until"`
		timeDelta int       `json:"time_delta"`
	}
	Result []cloudflare.ZoneAnalyticsColocation
}
