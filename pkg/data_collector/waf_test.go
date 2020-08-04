package data_collector

import (
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func Test_correctAdapting(t *testing.T) {
	aggregate := &model.Aggregate{
		ZoneName: "name",
		ZoneID:   "123",
		Totals:   make(map[time.Time]*model.Counters, 0),
	}

	GetWafTotals(aggregate, mockRequest(filepath.Join("..", "..", "test", "cloudlfare_graphql_response.json")))

	key, _ := time.Parse(time.RFC3339, "2020-07-26T10:01:00Z")
	assert.Equal(t, aggregate.Totals[key].WafTrigger["www.fr.lastminute.com"].Challenge.Key, "total.waf.trigger.challenge")
	assert.Equal(t, aggregate.Totals[key].WafTrigger["www.fr.lastminute.com"].Challenge.Value, 1)
}
