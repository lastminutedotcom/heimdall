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

func adaptDataToMetrics(datas []*model.Aggregate) []graphite.Metric {
	metrics := make([]graphite.Metric, 0)
	for _, data := range datas {
		metrics = append(metrics, metric(data.ZoneName, data.TotalRequestAll.Key, strconv.Itoa(data.TotalRequestAll.Value), data.Date))
		metrics = append(metrics, metric(data.ZoneName, data.TotalRequestCached.Key, strconv.Itoa(data.TotalRequestCached.Value), data.Date))
		metrics = append(metrics, metric(data.ZoneName, data.TotalRequestUncached.Key, strconv.Itoa(data.TotalRequestUncached.Value), data.Date))
		metrics = append(metrics, metric(data.ZoneName, data.TotalBandwidthAll.Key, strconv.Itoa(data.TotalBandwidthAll.Value), data.Date))
		metrics = append(metrics, metric(data.ZoneName, data.TotalBandwidthCached.Key, strconv.Itoa(data.TotalBandwidthCached.Value), data.Date))
		metrics = append(metrics, metric(data.ZoneName, data.TotalBandwidthUncached.Key, strconv.Itoa(data.TotalBandwidthUncached.Value), data.Date))

		for _, entry := range data.HTTPStatus {
			metrics = append(metrics, metric(data.ZoneName, entry.Key, strconv.Itoa(entry.Value), data.Date))
		}
	}
	return metrics
}

func metric(zone, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(METRICS_PREFIX + strings.Replace(zone, ".", "_", -1) + "." + key)

	logger.Printf("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}
