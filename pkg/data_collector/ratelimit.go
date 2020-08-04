package data_collector

import (
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"strings"
	"time"
)

func GetRatelimitTotals(aggregate *model.Aggregate, response *model.Response) (*model.Aggregate, error) {
	log.Info("Calculating rate limits for %s", aggregate.ZoneName)
	collectRateLimits(aggregate, response)
	return aggregate, nil
}

func collectRateLimits(aggregate *model.Aggregate, response *model.Response) {
	for _, zone := range response.Data.Viewer.Zones {
		for _, firewallEvent := range zone.FirewallEventsGroups {
			if !strings.EqualFold("rateLimit", firewallEvent.Dimensions.Source) {
				continue
			}
			OccurredAt := firewallEvent.Dimensions.OccurredAt.In(time.UTC)
			occurredAt := time.Date(OccurredAt.Year(), OccurredAt.Month(), OccurredAt.Day(), OccurredAt.Hour(), OccurredAt.Minute(), 0, 0, OccurredAt.Location())
			counters, exist := aggregate.Totals[occurredAt]
			if !exist {
				counters = model.NewCounters()
				aggregate.Totals[occurredAt] = counters
			}

			counter, present := counters.RateLimit[firewallEvent.Dimensions.Host]
			if !present {
				counter = model.NewRateLimitResult()
				counters.RateLimit[firewallEvent.Dimensions.Host] = counter
			}

			rateLimitCounters, _present := counter[firewallEvent.Dimensions.Method]

			if !_present {
				rateLimitCounters = model.NewSecurityEventCounters(strings.ToLower(firewallEvent.Dimensions.Method))
				counter[firewallEvent.Dimensions.Method] = rateLimitCounters
			}

			if firewallEvent.Dimensions.Action == "drop" {
				rateLimitCounters.Drop.Value++
			}
			if firewallEvent.Dimensions.Action == "simulate" {
				rateLimitCounters.Simulate.Value++
			}
			if firewallEvent.Dimensions.Action == "challenge" {
				rateLimitCounters.Challenge.Value++
			}
			if firewallEvent.Dimensions.Action == "jschallenge" {
				rateLimitCounters.JSChallenge.Value++
			}
			if firewallEvent.Dimensions.Action == "connectionClose" {
				rateLimitCounters.ConnectionClose.Value++
			}
		}
	}
}
