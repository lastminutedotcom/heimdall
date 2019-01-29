package metric

import (
	"fmt"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"github.com/marpaia/graphite-golang"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

const (
	wafMetricsPattern     = "cloudflare.%s.%s.%s"
	defaultMetricsPattern = "cloudflare.%s.%s"
)

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

			for host, entry := range counters.WafTrigger {
				metrics = append(metrics, wafMetric(aggregate.ZoneName, host, entry.Challenge.Key, strconv.Itoa(entry.Challenge.Value), date))
				metrics = append(metrics, wafMetric(aggregate.ZoneName, host, entry.JSChallenge.Key, strconv.Itoa(entry.JSChallenge.Value), date))
				metrics = append(metrics, wafMetric(aggregate.ZoneName, host, entry.Block.Key, strconv.Itoa(entry.Block.Value), date))
				metrics = append(metrics, wafMetric(aggregate.ZoneName, host, entry.Simulate.Key, strconv.Itoa(entry.Simulate.Value), date))
			}
		}
	}
	return metrics
}

func wafMetric(zone, host, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(fmt.Sprintf(wafMetricsPattern, normalize(zone), normalize(host), key))

	logger.Printf("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}

func metric(zone, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(fmt.Sprintf(defaultMetricsPattern, normalize(zone), key))

	logger.Printf("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}

func normalize(zone string) string {
	result := strings.Replace(zone, ".", "_", -1)
	return strings.Replace(result, " ", "_", -1)
}
