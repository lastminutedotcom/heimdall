package metric

import (
	"fmt"
	"git01.bravofly.com/N7/heimdall/pkg/logging"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"github.com/marpaia/graphite-golang"
	"strconv"
	"strings"
	"time"
)

const (
	hostMetricsPattern    = "cloudflare.%s.%s.%s"
	defaultMetricsPattern = "cloudflare.%s.%s"
)

func AdaptDataToMetrics(aggregates []*model.Aggregate) []graphite.Metric {
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
				metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.Challenge.Key, strconv.Itoa(entry.Challenge.Value), date))
				metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.JSChallenge.Key, strconv.Itoa(entry.JSChallenge.Value), date))
				metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.Block.Key, strconv.Itoa(entry.Block.Value), date))
				metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.Simulate.Key, strconv.Itoa(entry.Simulate.Value), date))
			}

			for host, rateLimitsCounters := range counters.RateLimit {
				for _, entry := range rateLimitsCounters {

					metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.Challenge.Key, strconv.Itoa(entry.Challenge.Value), date))
					metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.JSChallenge.Key, strconv.Itoa(entry.JSChallenge.Value), date))
					metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.ConnectionClose.Key, strconv.Itoa(entry.ConnectionClose.Value), date))
					metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.Drop.Key, strconv.Itoa(entry.Drop.Value), date))
					metrics = append(metrics, hostMetric(aggregate.ZoneName, host, entry.Simulate.Key, strconv.Itoa(entry.Simulate.Value), date))
				}

			}
		}
	}

	log.Info("added %d metrics", len(metrics))
	return metrics
}

func hostMetric(zone, host, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(fmt.Sprintf(hostMetricsPattern, normalize(zone), normalize(host), key))

	//log.Info("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}

func metric(zone, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(fmt.Sprintf(defaultMetricsPattern, normalize(zone), key))

	//log.Info("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}

func normalize(zone string) string {
	result := strings.Replace(zone, ".", "_", -1)
	return strings.Replace(result, " ", "_", -1)
}
