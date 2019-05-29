package ratelimit

import (
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"net/http"
	"time"
)

type RateLimitClient interface {
	GetSecurityEvents(zoneID string, since, until time.Time) ([]model.RateLimit, error)

	callSecurityEvent(url string) (*http.Response, model.RateLimitResponse, error)

	nextSecurityEventsBy(limits []model.RateLimit, result []model.RateLimit, zoneID string, nextPageId string, since, until time.Time) []model.RateLimit

	getSecurityEvent(zoneID string, nextPageId string) ([]model.RateLimit, string, error)
}
