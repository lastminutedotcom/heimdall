package zone

import (
	"git01.bravofly.com/n7/heimdall/pkg/client"
	"github.com/cloudflare/cloudflare-go"
)

type HttpZones struct {
}

func (HttpZones) GetZonesId() ([]cloudflare.Zone, error) {
	return client.CloudflareClient().ListZones()
}
