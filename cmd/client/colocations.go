package client

import (
	"github.com/cloudflare/cloudflare-go"
)

type ColocationsClient interface {
	GetColosAPI(zoneID string) ([]cloudflare.ZoneAnalyticsColocation, error)
}
