package data_collector

import (
	"git01.bravofly.com/n7/heimdall/cmd/client/colocation"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func GetColocationTotals(aggregates []*model.Aggregate, colocationsClient colocation.ColocationsClient) ([]*model.Aggregate, error) {
	for _, aggregate := range aggregates {
		logger.Printf("collecting co-location metrics for %s", aggregate.ZoneName)

		zoneAnalyticsDataArray, err := colocationsClient.GetColosAPI(aggregate.ZoneID)
		if err != nil {
			logger.Printf("ERROR Getting Analytics for zone %v, %v", aggregate.ZoneName, err)
			continue
		}

		collectColocation(zoneAnalyticsDataArray, aggregate)
	}
	return aggregates, nil
}

func collectColocation(zoneAnalyticsDataArray []cloudflare.ZoneAnalyticsColocation, aggregate *model.Aggregate) {
	for _, zoneAnalyticsData := range zoneAnalyticsDataArray {
		for _, timeSeries := range zoneAnalyticsData.Timeseries {

			counters, present := aggregate.Totals[timeSeries.Until]
			if !present {
				counters = model.NewCounters()
				aggregate.Totals[timeSeries.Until] = counters
			}

			counters.RequestAll.Value += timeSeries.Requests.All
			counters.RequestCached.Value += timeSeries.Requests.Cached
			counters.RequestUncached.Value += timeSeries.Requests.Uncached
			counters.BandwidthAll.Value += timeSeries.Bandwidth.All
			counters.BandwidthCached.Value += timeSeries.Bandwidth.Cached
			counters.BandwidthUncached.Value += timeSeries.Bandwidth.Uncached
			counters.HTTPStatus = totals(timeSeries.Requests.HTTPStatus, counters.HTTPStatus)

		}
	}
}

func totals(source map[string]int, target map[string]model.Counter) map[string]model.Counter {
	for k, v := range source {
		value := target[getKey(k)]
		value.Value += v
		target[getKey(k)] = value
	}
	return target
}

func getKey(httpCode string) string {
	if strings.HasPrefix(httpCode, "2") {
		return "2xx"
	}
	if strings.HasPrefix(httpCode, "3") {
		return "3xx"
	}
	if strings.HasPrefix(httpCode, "4") {
		return "4xx"
	}
	if strings.HasPrefix(httpCode, "5") {
		return "5xx"
	}

	return "1xx"
}
