package cmd

import (
	"git01.bravofly.com/N7/heimdall/pkg/client/colocation"
	"git01.bravofly.com/N7/heimdall/pkg/client/ratelimit"
	"git01.bravofly.com/N7/heimdall/pkg/client/waf"
	"git01.bravofly.com/N7/heimdall/pkg/client/zone"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"github.com/magiconair/properties/assert"
	"path/filepath"
	"sort"
	"testing"
)

func Test_integrationTest(t *testing.T) {

	mockZones := zone.MockZones{
		Path: filepath.Join("..", "test", "IT", "cloudflare_zone.json"),
	}

	mockColocations := colocation.MockColocations{
		Path: filepath.Join("..", "test", "IT", "cloudflare_colocation.json"),
	}

	mockWafs := waf.MockWafs{
		Path: filepath.Join("..", "test", "IT", "cloudflare_waf.json"),
	}

	mockRateLimitClient := ratelimit.MockRateLimitClient{
		Path: filepath.Join("..", "test", "IT", "cloudflare_ratelimit.json"),
	}

	aggregate := dataCollector(&model.Config{}, mockZones, mockColocations, mockWafs, mockRateLimitClient)
	metrics := adaptToMetrics(aggregate)

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Name < metrics[j].Name
	})

	assert.Equal(t, len(metrics), 337)
	assert.Equal(t, metrics[15].Name, "cloudflare.play_at.secure_play_at.total.ratelimit.post.challenge", )
	assert.Equal(t, metrics[24].Name, "cloudflare.play_at.secure_play_at.total.ratelimit.put.simulate")
	assert.Equal(t, metrics[15].Value, "4")
	assert.Equal(t, metrics[24].Value, "1")
}
