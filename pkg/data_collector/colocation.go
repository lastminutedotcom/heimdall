package data_collector

import (
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"strconv"
	"strings"
	"time"
)

func GetColocationTotals(aggregate *model.Aggregate, response *model.Response) (*model.Aggregate, error) {
	log.Info("Calculating colocation for %s", aggregate.ZoneName)
	collectColocation(aggregate, response)
	return aggregate, nil
}

func collectColocation(aggregate *model.Aggregate, response *model.Response) {
	for _, zone := range response.Data.Viewer.Zones {
		for _, group := range zone.HttpRequests1mGroups {
			key := time.Date(group.HttpRequestDimensions.DatetimeMinute.Year(), group.HttpRequestDimensions.DatetimeMinute.Month(),
				group.HttpRequestDimensions.DatetimeMinute.Day(), group.HttpRequestDimensions.DatetimeMinute.Hour(),
				group.HttpRequestDimensions.DatetimeMinute.Minute(), 0, 0, group.HttpRequestDimensions.DatetimeMinute.Location())

			counters, present := aggregate.Totals[key]
			if !present {
				counters = model.NewCounters()
				aggregate.Totals[key] = counters
			}
			counters.RequestAll.Value += group.HttpRequestSum.Requests
			counters.RequestCached.Value += group.HttpRequestSum.CachedRequests
			counters.RequestUncached.Value += group.HttpRequestSum.Requests - group.HttpRequestSum.CachedRequests
			counters.BandwidthAll.Value += group.HttpRequestSum.Bytes
			counters.BandwidthCached.Value += group.HttpRequestSum.CachedBytes
			counters.BandwidthUncached.Value += group.HttpRequestSum.Bytes - group.HttpRequestSum.CachedBytes
			counters.HTTPStatus = totals(group.HttpRequestSum.ResponseStatusMap, counters.HTTPStatus)
			aggregate.Totals[key] = counters
		}
	}
}

func totals(responseMap []model.ResponseStatusMap, target map[string]model.Counter) map[string]model.Counter {
	for _, a := range responseMap {
		value := target[getKey(strconv.Itoa(a.ResponseStatus))]
		value.Value += a.RequestCount
		target[getKey(strconv.Itoa(a.ResponseStatus))] = value
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
