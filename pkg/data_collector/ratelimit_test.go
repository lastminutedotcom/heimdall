package data_collector

import (
	"git01.bravofly.com/n7/heimdall/pkg/client/ratelimit"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_rateLimitDataCollection(t *testing.T) {

	aggregate := model.NewAggregate(cloudflare.Zone{
		ID:   "123",
		Name: "zone",
	})

	aggregates := make([]*model.Aggregate, 0)
	aggregates = append(aggregates, aggregate)

	GetRatelimitTotals(aggregates, &model.Config{}, ratelimit.MockRateLimitClient{})

	key, _ := time.Parse(time.RFC3339, "2019-01-17T13:22:00Z")
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["GET"].Challenge.Value, 0)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["POST"].Simulate.Value, 1)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["PATCH"].Simulate.Value, 1)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["PUT"].Simulate.Value, 1)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["DELETE"].Simulate.Value, 1)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["POST"].JSChallenge.Value, 2)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["POST"].ConnectionClose.Value, 1)
	assert.Equal(t, aggregate.Totals[key].RateLimit["secure.bravofly.at"]["POST"].Challenge.Value, 4)
}
