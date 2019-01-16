package model

import (
	"github.com/cloudflare/cloudflare-go"
	"time"
)

type Aggregate struct {
	ZoneName string
	ZoneID   string
	Date     time.Time

	TotalRequestAll        int
	TotalRequestCached     int
	TotalRequestUncached   int
	HTTPStatus             map[string]int
	TotalBandwidthAll      int
	TotalBandwidthCached   int
	TotalBandwidthUncached int
	TotalUniquesAll        int
	//WafTrigger       map[string]int
	//RateLimitTrigger map[string]int
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
