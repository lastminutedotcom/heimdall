package metric

import (
	"git01.bravofly.com/n7/heimdall/src/model"
	"github.com/marpaia/graphite-golang"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

const METRICS_PREFIX = "cloudflare.new."

func adaptDataToMetrics(aggregates []*model.Aggregate) []graphite.Metric {
	metrics := make([]graphite.Metric, 0)
	for _, aggregate := range aggregates {
		metrics = append(metrics, metric(aggregate.ZoneName, aggregate.TotalRequestAll.Key, strconv.Itoa(aggregate.TotalRequestAll.Value), aggregate.Date))
		metrics = append(metrics, metric(aggregate.ZoneName, aggregate.TotalRequestCached.Key, strconv.Itoa(aggregate.TotalRequestCached.Value), aggregate.Date))
		metrics = append(metrics, metric(aggregate.ZoneName, aggregate.TotalRequestUncached.Key, strconv.Itoa(aggregate.TotalRequestUncached.Value), aggregate.Date))
		metrics = append(metrics, metric(aggregate.ZoneName, aggregate.TotalBandwidthAll.Key, strconv.Itoa(aggregate.TotalBandwidthAll.Value), aggregate.Date))
		metrics = append(metrics, metric(aggregate.ZoneName, aggregate.TotalBandwidthCached.Key, strconv.Itoa(aggregate.TotalBandwidthCached.Value), aggregate.Date))
		metrics = append(metrics, metric(aggregate.ZoneName, aggregate.TotalBandwidthUncached.Key, strconv.Itoa(aggregate.TotalBandwidthUncached.Value), aggregate.Date))

		for _, entry := range aggregate.HTTPStatus {
			metrics = append(metrics, metric(aggregate.ZoneName, entry.Key, strconv.Itoa(entry.Value), aggregate.Date))
		}
	}
	return metrics
}

func metric(zone, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(METRICS_PREFIX + normalize(zone) + "." + key)

	logger.Printf("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}

func normalize(zone string) string {
	result := strings.Replace(zone, ".", "_", -1)
	return strings.Replace(result, " ", "_", -1)
}
