package client

import (
	"github.com/cloudflare/cloudflare-go"
)

func GetZonesId() ([]cloudflare.Zone, error) {
	return CloudflareClient().ListZones()
}
