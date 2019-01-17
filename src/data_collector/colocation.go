package data_collector

import (
	"git01.bravofly.com/n7/heimdall/src/client"
	"git01.bravofly.com/n7/heimdall/src/model"
	"log"
	"os"
	"strings"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func GetColocationTotals(aggregates []*model.Aggregate) ([]*model.Aggregate, error) {
	for _, aggregate := range aggregates {
		logger.Printf("collecting co-location metrics for %s", aggregate.ZoneName)

		zoneAnalyticsDataArray, date, err := client.GetColosAPI(aggregate.ZoneID)
		if err != nil {
			logger.Printf("ERROR Getting ZoneName Analytics for zone %v, %v", aggregate.ZoneName, err)
			return nil, err
		}

		aggregate.Date = date
		for _, zoneAnalyticsData := range zoneAnalyticsDataArray {
			for _, timeSeries := range zoneAnalyticsData.Timeseries {
				aggregate.TotalRequestAll.Value += timeSeries.Requests.All
				aggregate.TotalRequestCached.Value += timeSeries.Requests.Cached
				aggregate.TotalRequestUncached.Value += timeSeries.Requests.Uncached
				aggregate.TotalBandwidthAll.Value += timeSeries.Bandwidth.All
				aggregate.TotalBandwidthCached.Value += timeSeries.Bandwidth.Cached
				aggregate.TotalBandwidthUncached.Value += timeSeries.Bandwidth.Uncached

				aggregate.HTTPStatus = totals(timeSeries.Requests.HTTPStatus, aggregate.HTTPStatus)
			}
		}
	}
	return aggregates, nil
}

func totals(source map[string]int, target map[string]model.KeyValue) map[string]model.KeyValue {
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
