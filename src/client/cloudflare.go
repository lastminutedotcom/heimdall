package client

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/src/model"
	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"io/ioutil"
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
var rateLimiter = rate.NewLimiter(rate.Limit(4), 1)

var client = &http.Client{
	Timeout: time.Duration(5 * time.Second),
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

func GetColosAPI(zoneID string) ([]cloudflare.ZoneAnalyticsColocation, time.Time, error) {
	url := fmt.Sprintf(CloudFlareAPIRoot+"zones/%s/analytics/colos?since=%s&until=%s&continuous=%s", zoneID, "-1", "-1", "false")
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := doHttpCall(request)
	if err != nil {
		return nil, time.Now(), fmt.Errorf("get colocation analytics HTTP call error: %v", err)
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

func GetZonesId() ([]*model.Aggregate, error) {
	zones, err := cloudflareClient().ListZones()
	if err != nil {
		logger.Printf("ERROR ZoneName from CF Client %v", zones)
		return nil, err
	}

	result := make([]*model.Aggregate, 0)
	for _, zone := range zones {
		result = append(result, model.NewAggregate(zone))
	}

	return result, nil
}
