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
		metrics = append(metrics, metric(data.ZoneName, "total.requests.all", strconv.Itoa(data.TotalRequestAll), data.Date))
		metrics = append(metrics, metric(data.ZoneName, "total.requests.cached", strconv.Itoa(data.TotalRequestCached), data.Date))
		metrics = append(metrics, metric(data.ZoneName, "total.requests.uncached", strconv.Itoa(data.TotalRequestUncached), data.Date))
		metrics = append(metrics, metric(data.ZoneName, "total.bandwidth.all", strconv.Itoa(data.TotalBandwidthAll), data.Date))
		metrics = append(metrics, metric(data.ZoneName, "total.bandwidth.cached", strconv.Itoa(data.TotalBandwidthCached), data.Date))
		metrics = append(metrics, metric(data.ZoneName, "total.bandwidth.uncached", strconv.Itoa(data.TotalBandwidthUncached), data.Date))

		for httpFamily, counter := range data.HTTPStatus {
			metrics = append(metrics, metric(data.ZoneName, "total.requests.http_status."+httpFamily, strconv.Itoa(counter), data.Date))
		}
	}
	return metrics
}

func metric(zone, key, value string, date time.Time) graphite.Metric {
	metricKey := strings.ToLower(METRICS_PREFIX + strings.Replace(zone, ".", "_", -1) + "." + key)

	logger.Printf("added metric %s, value %s, %v", metricKey, value, date.Unix())

	return graphite.NewMetric(metricKey, value, date.Unix())
}
