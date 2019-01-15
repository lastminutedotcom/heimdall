package main

import "C"
import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/marpaia/graphite-golang"
	"gopkg.in/robfig/cron.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

const (
	CloudFlareAPIRoot = "https://api.cloudflare.com/client/v4/"
	key               = "f73d2fd09a50dd1234a26d37e794de982fc0c"
	email             = "api.sre@lastminute.com"
	orgId             = "f5fd3b3741817e2080883b52b5995643"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

type DataAggregation struct {
	ZoneName         string
	ZoneID           string
	Totals           *cloudflare.ZoneAnalytics
	WafTrigger       map[string]int
	RateLimitTrigger map[string]int
}

var client = &http.Client{
	Timeout: time.Duration(30 * time.Second),
}

func main() {

	cronExp := "0 * * * * *"
	logger.Printf("start collecting data %s", cronExp)

	c := cron.New()
	c.AddFunc(cronExp, collectingData)

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig

	c.Stop()
	fmt.Println("Got signal:", s)

}

func collectingData() {
	aggregations, _ := getZonesId(cloudflareClient())
	aggregations, _ = getColocationTotals(cloudflareClient(), aggregations)

	//for _, aggregation := range aggregations {
	//	logger.Printf("------------------------------------------")
	//	logger.Printf("aggregation name %s", aggregation.ZoneName)
	//	logger.Printf("aggregation id %s", aggregation.ZoneID)
	//	logger.Printf("request http status %v", aggregation.Totals.Requests.HTTPStatus)
	//	logger.Printf("request http al %v", aggregation.Totals.Requests.All)
	//	logger.Printf("request http cached %v", aggregation.Totals.Requests.Cached)
	//	logger.Printf("request http uncached %v", aggregation.Totals.Requests.Uncached)
	//	logger.Printf("bandwidth all %v", aggregation.Totals.Bandwidth.All)
	//	logger.Printf("bandwidth cached %v", aggregation.Totals.Bandwidth.Cached)
	//	logger.Printf("bandwidth uncached %v", aggregation.Totals.Bandwidth.Uncached)
	//	logger.Printf("uniques all %v", aggregation.Totals.Uniques.All)
	//
	//	logger.Printf("------------------------------------------")
	//
	//}

	pushMetrics(aggregations)
}

func getZonesId(client *cloudflare.API) ([]*DataAggregation, error) {
	zones, err := client.ListZones()
	if err != nil {
		logger.Printf("ERROR ZoneName from CF Client %v", zones)
		return nil, err
	}

	result := make([]*DataAggregation, 0)
	for _, zone := range zones {
		result = append(result, &DataAggregation{ZoneName: zone.Name, ZoneID: zone.ID})
	}

	return result, nil
}

func getColocationTotals(client *cloudflare.API, dataAggregations []*DataAggregation) ([]*DataAggregation, error) {

	// doesn't work with values minor then -60
	now := time.Now()
	since := now.Add(time.Duration(-60) * time.Minute)
	continuous := false
	logger.Printf("data from %s to %s", since.Format(time.RFC3339), now.Format(time.RFC3339))

	options := cloudflare.ZoneAnalyticsOptions{
		Until:      &now,
		Since:      &since,
		Continuous: &continuous,
	}

	for _, data := range dataAggregations {
		logger.Printf("collecting metrics for %s", data.ZoneName)
		zoneAnalyticsData, err := client.ZoneAnalyticsDashboard(data.ZoneID,
			options)
		if err != nil {
			logger.Printf("ERROR Getting ZoneName Analytics for zone %v, %v", data.ZoneName, err)
			return nil, err
		}
		data.Totals = &zoneAnalyticsData.Totals
	}
	return dataAggregations, nil
}

func cloudflareClient() *cloudflare.API {
	c, err := cloudflare.New(key, email, cloudflare.UsingOrganization(orgId), cloudflare.UsingRateLimit(2))
	if err != nil {
		logger.Fatalf("could not create client for Cloudflare: %v", err)
	}
	return c
}

func doHttpCall(request *http.Request) (*http.Response, error) {
	request = setHeaders(request)
	return client.Do(request)
}

func setHeaders(request *http.Request) *http.Request {
	for key, value := range createHeaders() {
		request.Header.Set(key, value)
	}
	return request
}

func createHeaders() map[string]string {
	return map[string]string{
		"X-Auth-Email": email,
		"X-Auth-Key":   key,
		"Content-Type": "application/json",
	}
}

func pushMetrics(datas []*DataAggregation) {

	newGraphite, err := graphite.NewGraphite("10.120.172.134", 2113)

	if err != nil {
		newGraphite = graphite.NewGraphiteNop("10.120.172.134", 2113)
	}

	metrics := make([]graphite.Metric, 0)
	for _, data := range datas {
		metrics = append(metrics, metric(data.ZoneName, "total.requests.all", strconv.Itoa(data.Totals.Requests.All)))
		metrics = append(metrics, metric(data.ZoneName, "total.requests.cached", strconv.Itoa(data.Totals.Requests.Cached)))
		metrics = append(metrics, metric(data.ZoneName, "total.requests.uncached", strconv.Itoa(data.Totals.Requests.Uncached)))
	}
	newGraphite.SendMetrics(metrics)
}

func metric(zone, key, value string) graphite.Metric {
	metricKey := strings.ToLower("cloudflare.new." + strings.Replace(zone, ".", "_", -1) + "." + key)

	logger.Printf("added metric %s, value %s, %v", metricKey, value, time.Now().Unix())

	return graphite.NewMetric(metricKey, value, time.Now().Unix())
}
