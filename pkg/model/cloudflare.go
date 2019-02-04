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

type WAFResponse struct {
	cloudflare.Response
	ResultInfo  ResultInfo   `json:"result_info"`
	WafTriggers []WafTrigger `json:"result"`
}

type ResultInfo struct {
	NextPageId string `json:"next_page_id"`
}

type WafTrigger struct {
	Host       string    `json:"host"`
	Action     string    `json:"action"`
	OccurredAt time.Time `json:"occurred_at"`
}

type RateLimitResponse struct {
	cloudflare.Response
	ResultInfo struct {
		Cursors struct {
			After  string `json:"after"`
			Before string `json:"before"`
		}
	}
	Result []RateLimit `json:"result"`
}

type RateLimit struct {
	Host       string    `json:"host"`
	Action     string    `json:"action"`
	Method     string    `json:"method"`
	OccurredAt time.Time `json:"occurred_at"`
}
