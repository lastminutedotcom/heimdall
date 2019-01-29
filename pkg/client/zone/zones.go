package zone

import "github.com/cloudflare/cloudflare-go"

type ZonesClient interface {
	GetZonesId() ([]cloudflare.Zone, error)
}
