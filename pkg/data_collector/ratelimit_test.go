package data_collector

import (
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/magiconair/properties/assert"
	"path/filepath"
	"testing"
	"time"
)

func Test_ratelimitDataCollection(t *testing.T) {
	aggregate := &model.Aggregate{
		ZoneName: "name",
		ZoneID:   "123",
		Totals:   make(map[time.Time]*model.Counters, 0),
	}

	GetRatelimitTotals(aggregate, mockRequest(filepath.Join("..", "..", "test", "cloudlfare_graphql_response.json")))

	key, _ := time.Parse(time.RFC3339, "2020-07-26T10:01:00Z")
	assert.Equal(t, aggregate.Totals[key].RateLimit["www.fr.lastminute.com"]["GET"].Drop.Key, "total.ratelimit.get.drop")
	assert.Equal(t, aggregate.Totals[key].RateLimit["www.fr.lastminute.com"]["GET"].Drop.Value, 1)
}
