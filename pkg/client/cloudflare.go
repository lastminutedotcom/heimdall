package client

import (
	"golang.org/x/net/context"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	CloudFlareAPIRoot = "https://api.cloudflare.com/client/v4/"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)
var rateLimiter = rate.NewLimiter(rate.Limit(3), 1) // 3rps (900 req/5 min)

var client = &http.Client{
	Timeout: time.Duration(20 * time.Second),
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
	}
}
