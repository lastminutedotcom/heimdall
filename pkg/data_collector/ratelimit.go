package data_collector

import (
	"git01.bravofly.com/n7/heimdall/pkg/client/ratelimit"
	"git01.bravofly.com/n7/heimdall/pkg/logging"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"strconv"
	"strings"
	"time"
)

func GetRatelimitTotals(aggregates []*model.Aggregate, config *model.Config, client ratelimit.RateLimitClient) ([]*model.Aggregate, error) {
	for _, aggregate := range aggregates {
		log.Info("collecting rate limit metrics for %s", aggregate.ZoneName)

		utc, _ := time.LoadLocation("UTC")
		until := time.Now().In(utc)
		everyMinutes, _ := strconv.Atoi(config.CollectEveryMinutes)
		since := until.Add(time.Duration(everyMinutes) * time.Minute * -1)

		rateLimits, err := client.GetSecurityEvents(aggregate.ZoneID, since, until)
		if err != nil {
			log.Error("ERROR Getting rate limit trigger for zone %v, %v", aggregate.ZoneName, err)
			continue
		}
		collectRateLimits(rateLimits, utc, aggregate)
	}

	return aggregates, nil
}

func collectRateLimits(rateLimits []model.RateLimit, utc *time.Location, aggregate *model.Aggregate) {
	for _, rateLimit := range rateLimits {
		OccurredAt := rateLimit.OccurredAt.In(utc)
		occurredAt := time.Date(OccurredAt.Year(), OccurredAt.Month(), OccurredAt.Day(), OccurredAt.Hour(), OccurredAt.Minute(), 0, 0, OccurredAt.Location())
		counters, exist := aggregate.Totals[occurredAt]
		if !exist {
			counters = model.NewCounters()
			aggregate.Totals[occurredAt] = counters
		}

		counter, present := counters.RateLimit[rateLimit.Host]
		if !present {
			counter = model.NewRateLimitResult()
			counters.RateLimit[rateLimit.Host] = counter
		}

		rateLimitCounters, _present := counter[rateLimit.Method]

		if !_present {
			rateLimitCounters = model.NewSecurityEventCounters(strings.ToLower(rateLimit.Method))
			counter[rateLimit.Method] = rateLimitCounters
		}

		if rateLimit.Action == "drop" {
			rateLimitCounters.Drop.Value++
		}
		if rateLimit.Action == "simulate" {
			rateLimitCounters.Simulate.Value++
		}
		if rateLimit.Action == "challenge" {
			rateLimitCounters.Challenge.Value++
		}
		if rateLimit.Action == "jschallenge" {
			rateLimitCounters.JSChallenge.Value++
		}
		if rateLimit.Action == "connectionClose" {
			rateLimitCounters.ConnectionClose.Value++
		}
	}
}
