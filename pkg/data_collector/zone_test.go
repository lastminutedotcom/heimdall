package data_collector

import (
	"github.com/lastminutedotcom/heimdall/pkg/client/zone"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_zoneCollection(t *testing.T) {
	aggregates, _ := GetZones(zone.MockZones{
		Path: filepath.Join("..", "..", "test", "cloudflare_zone.json"),
	})

	assert.Equal(t, len(aggregates), 11)
	assert.Equal(t, aggregates[0].ZoneName, "play.at")
	assert.Equal(t, aggregates[10].ZoneName, "lastplay.com")
	assert.Equal(t, aggregates[0].ZoneID, "11111111111111111111111111111111")
	assert.Equal(t, aggregates[10].ZoneID, "00000000000000000000000000000000")
}
