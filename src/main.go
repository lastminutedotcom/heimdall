package main

import "C"
import (
	"fmt"
	"git01.bravofly.com/n7/heimdall/src/client"
	"git01.bravofly.com/n7/heimdall/src/metric"
	"git01.bravofly.com/n7/heimdall/src/model"
	"gopkg.in/robfig/cron.v2"
	"log"
	"os"
	"os/signal"
	"strings"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func main() {
	cronExp := "0 * * * * *"
	//cronExp := "* * * * * *"
	logger.Printf("start collecting data %s", cronExp)

	c := cron.New()
	c.AddFunc(cronExp, orchestrator)

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig

	c.Stop()
	fmt.Println("Got signal:", s)

	//orchestrator()

}

func orchestrator() {
	aggregate := dataCollector()
	metric.PushMetrics(aggregate)
}

func dataCollector() []*model.Aggregate {
	aggregate, _ := client.GetZonesId()
	aggregate, _ = getColocationTotals(aggregate)
	return aggregate
}

func getColocationTotals(dataAggregations []*model.Aggregate) ([]*model.Aggregate, error) {
	for _, data := range dataAggregations {
		logger.Printf("collecting metrics for %s", data.ZoneName)

		zoneAnalyticsDataArray, date, err := client.GetColosAPI(data.ZoneID)
		if err != nil {
			logger.Printf("ERROR Getting ZoneName Analytics for zone %v, %v", data.ZoneName, err)
			return nil, err
		}

		data.Date = date
		for _, zoneAnalyticsData := range zoneAnalyticsDataArray {
			for _, timeSeries := range zoneAnalyticsData.Timeseries {
				data.TotalRequestAll.Value += timeSeries.Requests.All
				data.TotalRequestCached.Value += timeSeries.Requests.Cached
				data.TotalRequestUncached.Value += timeSeries.Requests.Uncached
				data.TotalBandwidthAll.Value += timeSeries.Bandwidth.All
				data.TotalBandwidthCached.Value += timeSeries.Bandwidth.Cached
				data.TotalBandwidthUncached.Value += timeSeries.Bandwidth.Uncached

				data.HTTPStatus = totals(timeSeries.Requests.HTTPStatus, data.HTTPStatus)
			}
		}
	}
	return dataAggregations, nil
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
