package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	CloudFlareAPIRoot        = "https://api.cloudflare.com/client/v4/"
	CloudFlareGraphQLAPIRoot = "https://api.cloudflare.com/client/v4/graphql"
)

var rateLimiter = rate.NewLimiter(rate.Limit(3), 1) // 3rps (900 req/5 min)

var client = &http.Client{
	Timeout: time.Duration(20 * time.Second),
}

func CallGraphQlApi(zoneId string, startDate, endDate time.Time) *model.Response {

	resp, err := DoHttpCall(createRequest(zoneId, startDate, endDate))
	if err != nil {
		log.Info("HTTP call error: %v", err)
	}
	response := &model.Response{}

	if err != nil {
		log.Info("HTTP call error: %v", err)
	}

	b, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(b, &response); err != nil {
		log.Info("HTTP body marshal to JSON error: %v", err)
	}

	return response
}

func createRequest(zoneId string, startDate, endDate time.Time) *http.Request {
	url := fmt.Sprintf(CloudFlareGraphQLAPIRoot)

	req := fmt.Sprintf("{viewer{zones(filter: {zoneTag: \"%s\"}) {httpRequests1mGroups(orderBy: [datetimeMinute_ASC], limit: 10000, filter: {datetime_gt: \"%s\", datetime_lt: \"%s\"}) {dimensions {datetimeMinute} sum { bytes cachedBytes cachedRequests requests responseStatusMap { requests edgeResponseStatus}}}firewallEventsGroups(limit: 10000, filter: {datetime_gt: \"%s\", datetime_lt: \"%s\"}) {dimensions { action occurredDatetime clientRequestHTTPHost clientRequestHTTPMethodName source}}}}}",
		zoneId, startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z"), startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z"))

	var requestBody bytes.Buffer
	requestBodyObj := struct {
		Query string `json:"query"`
	}{
		Query: req,
	}
	if err := json.NewEncoder(&requestBody).Encode(requestBodyObj); err != nil {
		log.Info("Encoding error: %v", err)
	}

	request, _ := http.NewRequest(http.MethodPost, url, &requestBody)
	return request
}

func DoHttpCall(request *http.Request) (*http.Response, error) {
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
		"X-Auth-Email": os.Getenv("CLOUDFLARE_EMAIL"),
		"X-Auth-Key":   os.Getenv("CLOUDFLARE_TOKEN"),
		"Content-Type": "application/json",
		"Accept":       "*/*",
	}
}
