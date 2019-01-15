package main

import (
	"github.com/cloudflare/cloudflare-go"
	"github.com/marpaia/graphite-golang"
	"log"
	"net/http"
	"os"
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
	ColocationTotals *cloudflare.ZoneAnalytics
	WafTrigger       map[string]int
	RateLimitTrigger map[string]int
}

var client = &http.Client{
	Timeout: time.Duration(30 * time.Second),
}

func main() {
	collectingData()
}

func collectingData() {
	aggregations, _ := getZonesId(cloudflareClient())
	aggregations, _ = getColocationTotals(cloudflareClient(), aggregations)

	for _, aggregation := range aggregations {
		logger.Printf("------------------------------------------")
		logger.Printf("aggregation name %s", aggregation.ZoneName)
		logger.Printf("aggregation id %s", aggregation.ZoneID)
		logger.Printf("request http status %v", aggregation.ColocationTotals.Requests.HTTPStatus)
		logger.Printf("request http al %v", aggregation.ColocationTotals.Requests.All)
		logger.Printf("request http cached %v", aggregation.ColocationTotals.Requests.Cached)
		logger.Printf("request http uncached %v", aggregation.ColocationTotals.Requests.Uncached)
		logger.Printf("bandwidth all %v", aggregation.ColocationTotals.Bandwidth.All)
		logger.Printf("bandwidth cached %v", aggregation.ColocationTotals.Bandwidth.Cached)
		logger.Printf("bandwidth uncached %v", aggregation.ColocationTotals.Bandwidth.Uncached)
		logger.Printf("uniques all %v", aggregation.ColocationTotals.Uniques.All)

		logger.Printf("------------------------------------------")

	}

	//pushMetrics(aggregations)
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
		zoneAnalyticsData, err := client.ZoneAnalyticsDashboard(data.ZoneID,
			options)
		if err != nil {
			logger.Printf("ERROR Getting ZoneName Analytics for zone %v, %v", data.ZoneName, err)
			return nil, err
		}
		data.ColocationTotals = &zoneAnalyticsData.Totals
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

func pushMetrics(data []*DataAggregation) {

	newGraphite, err := graphite.NewGraphite("10.120.172.134", 2113)

	if err != nil {
		newGraphite = graphite.NewGraphiteNop("10.120.172.134", 2113)
	}

	metrics := make([]graphite.Metric, 1)
	metrics[0] = graphite.NewMetric("ColocationTotals.Requests.HTTPStatus", "1000", time.Now().Unix())
	newGraphite.SendMetrics(metrics)
}
