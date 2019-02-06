package data_collector

import (
	"git01.bravofly.com/N7/heimdall/pkg/client/waf"
	"git01.bravofly.com/N7/heimdall/pkg/logging"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"strconv"
	"time"
)

func GetWafTotals(aggregates []*model.Aggregate, config *model.Config, client waf.WafsClient) ([]*model.Aggregate, error) {
	for _, aggregate := range aggregates {
		log.Info("collecting waf trigger metrics for %s", aggregate.ZoneName)
		utc, _ := time.LoadLocation("UTC")
		since := time.Now().In(utc)
		everyMinutes, _ := strconv.Atoi(config.CollectEveryMinutes)
		until := since.Add(time.Duration(everyMinutes) * time.Minute * -1)

		triggers, err := client.GetWafTriggersBy(aggregate.ZoneID, since, until)
		if err != nil {
			log.Error("ERROR Getting WAF trigger for zone %v, %v", aggregate.ZoneName, err)
			continue
		}

		collectWaf(triggers, utc, aggregate)
	}

	return aggregates, nil

}

func collectWaf(triggers []model.WafTrigger, utc *time.Location, aggregate *model.Aggregate) {
	for _, trigger := range triggers {
		OccurredAt := trigger.OccurredAt.In(utc)
		occurredAt := time.Date(OccurredAt.Year(), OccurredAt.Month(), OccurredAt.Day(), OccurredAt.Hour(), OccurredAt.Minute(), 0, 0, OccurredAt.Location())
		counters, exist := aggregate.Totals[occurredAt]
		if !exist {
			counters = model.NewCounters()
			aggregate.Totals[occurredAt] = counters
		}

		counter, present := counters.WafTrigger[trigger.Host]
		if !present {
			counter = model.NewWafTriggerResult()
			counters.WafTrigger[trigger.Host] = counter
		}

		if trigger.Action == "block" {
			counter.Block.Value++
		}
		if trigger.Action == "challenge" {
			counter.Challenge.Value++
		}
		if trigger.Action == "jschallenge" {
			counter.JSChallenge.Value++
		}
		if trigger.Action == "simulate" {
			counter.Simulate.Value++
		}
	}
}
