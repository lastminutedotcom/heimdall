package data_collector

import (
	"git01.bravofly.com/N7/heimdall/pkg/client/colocation"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"github.com/cloudflare/cloudflare-go"
	"github.com/magiconair/properties/assert"
	"path/filepath"
	"testing"
	"time"
)

func Test_colocationDataCollection(t *testing.T) {
	aggregate := model.NewAggregate(cloudflare.Zone{
		ID:   "123",
		Name: "zone",
	})

	aggregates := make([]*model.Aggregate, 0)
	aggregates = append(aggregates, aggregate)

	GetColocationTotals(aggregates, colocation.MockColocations{
		Path: filepath.Join("..", "..", "test", "cloudflare_colocation.json"),
	})

	assert.Equal(t, len(aggregate.Totals), 5)
	key, _ := time.Parse(time.RFC3339, "2019-01-23T15:01:00Z")
	assert.Equal(t, aggregate.Totals[key].RequestAll.Value, 1033)
	assert.Equal(t, aggregate.Totals[key].BandwidthCached.Value, 2821897)
	assert.Equal(t, aggregate.Totals[key].HTTPStatus["2xx"].Value, 861)

}
