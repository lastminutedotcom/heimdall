package data_collector

import (
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"strings"
	"time"
)

func GetWafTotals(aggregate *model.Aggregate, response *model.Response) (*model.Aggregate, error) {
	collectWaf(aggregate, response)
	return aggregate, nil
}

func collectWaf(aggregate *model.Aggregate, response *model.Response) {
	for _, zone := range response.Data.Viewer.Zones {
		for _, firewallEvent := range zone.FirewallEventsGroups {
			if strings.EqualFold("rateLimit", firewallEvent.Dimensions.Source) {
				continue
			}

			OccurredAt := firewallEvent.Dimensions.OccurredAt.In(time.UTC)
			occurredAt := time.Date(OccurredAt.Year(), OccurredAt.Month(), OccurredAt.Day(), OccurredAt.Hour(), OccurredAt.Minute(), 0, 0, OccurredAt.Location())
			counters, exist := aggregate.Totals[occurredAt]

			if !exist {
				counters = model.NewCounters()
				aggregate.Totals[occurredAt] = counters
			}

			counter, present := counters.WafTrigger[firewallEvent.Dimensions.Host]
			if !present {
				counter = model.NewWafTriggerResult()
				counters.WafTrigger[firewallEvent.Dimensions.Host] = counter
			}

			if firewallEvent.Dimensions.Action == "block" {
				counter.Block.Value++
			}
			if firewallEvent.Dimensions.Action == "challenge" {
				counter.Challenge.Value++
			}
			if firewallEvent.Dimensions.Action == "jschallenge" {
				counter.JSChallenge.Value++
			}
			if firewallEvent.Dimensions.Action == "simulate" {
				counter.Simulate.Value++
			}
		}
	}
}
