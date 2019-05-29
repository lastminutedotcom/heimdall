package data_collector

import (
	"git01.bravofly.com/N7/heimdall/pkg/client/zone"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_zoneCollection(t *testing.T) {
	aggregates, _ := GetZones(zone.MockZones{
		Path: filepath.Join("..", "..", "test", "cloudflare_zone.json"),
	})

	assert.Equal(t, len(aggregates), 18)
	assert.Equal(t, aggregates[0].ZoneName, "play.at")
	assert.Equal(t, aggregates[12].ZoneName, "jumbo.com")
	assert.Equal(t, aggregates[0].ZoneID, "aaaaaaaaaaaabbbbbbbbbbbbbbbbbbcc")
	assert.Equal(t, aggregates[12].ZoneID, "eeeeeeeeeeeeebbbbbbbbbbbbbbbbbbc")
}
