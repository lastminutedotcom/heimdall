package data_collector

import (
	"git01.bravofly.com/n7/heimdall/cmd/client"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"github.com/cloudflare/cloudflare-go"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func readAnalyticsColocationResponse() []cloudflare.ZoneAnalyticsColocation {
	colocations, _ := client.MockColocations{}.GetColosAPI("123")
	return colocations
}

func Test_colocationDataCollection(t *testing.T) {
	aggregate := model.NewAggregate(cloudflare.Zone{
		ID:   "123",
		Name: "zone",
	})

	collectColocation(readAnalyticsColocationResponse(), aggregate)

	assert.Equal(t, len(aggregate.Totals), 5)
	key, _ := time.Parse(time.RFC3339, "2019-01-23T15:01:00Z")
	assert.Equal(t, aggregate.Totals[key].RequestAll.Value, 1033)
	assert.Equal(t, aggregate.Totals[key].BandwidthCached.Value, 2821897)
	assert.Equal(t, aggregate.Totals[key].HTTPStatus["2xx"].Value, 861)

}
