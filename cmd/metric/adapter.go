package metric

import (
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"github.com/marpaia/graphite-golang"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

const METRICS_PREFIX = "cloudflare."

func adaptDataToMetrics(aggregates []*model.Aggregate) []graphite.Metric {
	metrics := make([]graphite.Metric, 0)
	for _, aggregate := range aggregates {
		for date, counters := range aggregate.Totals {
			metrics = append(metrics, metric(aggregate.ZoneName, counters.RequestAll.Key, strconv.Itoa(counters.RequestAll.Value), date))
			metrics = append(metrics, metric(aggregate.ZoneName, counters.RequestCached.Key, strconv.Itoa(counters.RequestCached.Value), date))
			metrics = append(metrics, metric(aggregate.ZoneName, counters.RequestUncached.Key, strconv.Itoa(counters.RequestUncached.Value), date))
			metrics = append(metrics, metric(aggregate.ZoneName, counters.BandwidthAll.Key, strconv.Itoa(counters.BandwidthAll.Value), date))
			metrics = append(metrics, metric(aggregate.ZoneName, counters.BandwidthCached.Key, strconv.Itoa(counters.BandwidthCached.Value), date))
			metrics = append(metrics, metric(aggregate.ZoneName, counters.BandwidthUncached.Key, strconv.Itoa(counters.BandwidthUncached.Value), date))

			for _, entry := range counters.HTTPStatus {
				metrics = append(metrics, metric(aggregate.ZoneName, entry.Key, strconv.Itoa(entry.Value), date))
			}
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
