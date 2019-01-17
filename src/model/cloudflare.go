package model

import (
	"github.com/cloudflare/cloudflare-go"
	"time"
)

type ZoneAnalyticsColocationResponse struct {
	cloudflare.Response
	Query struct {
		Since     time.Time `json:"since"`
		Until     time.Time `json:"until"`
		timeDelta int       `json:"time_delta"`
	}
	Result []cloudflare.ZoneAnalyticsColocation
}
