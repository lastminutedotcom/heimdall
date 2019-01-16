package main

import "C"
import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/src/metric"
	"git01.bravofly.com/n7/heimdall/src/model"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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
var rateLimiter = rate.NewLimiter(rate.Limit(4), 1)

var client = &http.Client{
	Timeout: time.Duration(5 * time.Second),
}

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
	aggregate, _ := getZonesId(cloudflareClient())
	aggregate, _ = getColocationTotals(aggregate)
	return aggregate
}

func getZonesId(client *cloudflare.API) ([]*model.Aggregate, error) {
	zones, err := client.ListZones()
	if err != nil {
		logger.Printf("ERROR ZoneName from CF Client %v", zones)
		return nil, err
	}

	result := make([]*model.Aggregate, 0)
	for _, zone := range zones {
		result = append(result, &model.Aggregate{
			ZoneName:               zone.Name,
			ZoneID:                 zone.ID,
			TotalRequestAll:        model.KeyValue{Key: "total.requests.all", Value: 0},
			TotalRequestCached:     model.KeyValue{Key: "total.requests.cached", Value: 0},
			TotalRequestUncached:   model.KeyValue{Key: "total.requests.uncached", Value: 0},
			TotalBandwidthAll:      model.KeyValue{Key: "total.bandwidth.all", Value: 0},
			TotalBandwidthCached:   model.KeyValue{Key: "total.bandwidth.cached", Value: 0},
			TotalBandwidthUncached: model.KeyValue{Key: "total.bandwidth.uncached", Value: 0},
			HTTPStatus:             map[string]int{"2xx": 0, "3xx": 0, "4xx": 0, "5xx": 0},
		})
	}

	return result, nil
}

func getColocationTotals(dataAggregations []*model.Aggregate) ([]*model.Aggregate, error) {
	for _, data := range dataAggregations {
		logger.Printf("collecting metrics for %s", data.ZoneName)

		zoneAnalyticsDataArray, date, err := callColocationAnalyticsAPI(data.ZoneID)
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

func totals(source, target map[string]int) map[string]int {

	for k, v := range source {
		key := getKey(k)
		if value, present := target[key]; present {
			value += v
			target[key] = value
		} else {
			target[key] = v
		}
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

func callColocationAnalyticsAPI(zoneID string) ([]cloudflare.ZoneAnalyticsColocation, time.Time, error) {
	url := fmt.Sprintf(CloudFlareAPIRoot+"zones/%s/analytics/colos?since=%s&until=%s&continuous=%s", zoneID, "-1", "-1", "false")
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := doHttpCall(request)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("get zones HTTP call error: %v", err)
	}
	response := model.ZoneAnalyticsColocationResponse{}
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, time.Now(), fmt.Errorf("HTTP body marshal to JSON error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		return response.Result, response.Query.Until, nil
	}
	return nil, time.Now(), fmt.Errorf("get colocation analytics HTTP error %d", resp.StatusCode)
}

func cloudflareClient() *cloudflare.API {
	c, err := cloudflare.New(key, email, cloudflare.UsingOrganization(orgId), cloudflare.UsingRateLimit(2), cloudflare.HTTPClient(client))
	if err != nil {
		logger.Fatalf("could not create client for Cloudflare: %v", err)
	}
	return c
}

func doHttpCall(request *http.Request) (*http.Response, error) {
	rateLimiter.Wait(context.TODO())
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
