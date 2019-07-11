package data_collector

import (
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/cloudflare/cloudflare-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_correctAdapting(t *testing.T) {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.UTC)

	triggers := make([]model.WafTrigger, 0)
	triggers = append(triggers, newWafTrigger("host.it", "block", now))
	triggers = append(triggers, newWafTrigger("host.it", "block", now))
	triggers = append(triggers, newWafTrigger("host.it", "challenge", now))
	triggers = append(triggers, newWafTrigger("host.it", "challenge", now))
	triggers = append(triggers, newWafTrigger("host.it", "challenge", now))
	triggers = append(triggers, newWafTrigger("host.it", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.it", "jschallenge", now))

	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))
	triggers = append(triggers, newWafTrigger("host.com", "simulate", now))

	aggregate := model.NewAggregate(cloudflare.Zone{
		ID:   "123",
		Name: "zone",
	})
	aggregate.Totals[now] = model.NewCounters()

	collectWaf(triggers, aggregate)

	assert.Equal(t, aggregate.Totals[now].WafTrigger["host.com"].Simulate.Value, 7)
	assert.Equal(t, aggregate.Totals[now].WafTrigger["host.it"].Block.Value, 2)

}

func newWafTrigger(host, action string, now time.Time) model.WafTrigger {
	return model.WafTrigger{
		Host:       host,
		Action:     action,
		OccurredAt: now,
	}
}
