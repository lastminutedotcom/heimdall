package client

import (
	"github.com/cloudflare/cloudflare-go"
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
var rateLimiter = rate.NewLimiter(rate.Limit(4), 1)

var client = &http.Client{
	Timeout: time.Duration(10 * time.Second),
}

func cloudflareClient() *cloudflare.API {
	c, err := cloudflare.New(os.Getenv("CLOUDFLARE_TOKEN"), os.Getenv("CLOUDFLARE_EMAIL"),
		cloudflare.UsingOrganization(os.Getenv("CLOUDFLARE_ORG_ID")), cloudflare.HTTPClient(client))
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
		"X-Auth-Email": os.Getenv("CLOUDFLARE_EMAIL"),
		"X-Auth-Key":   os.Getenv("CLOUDFLARE_TOKEN"),
		"Content-Type": "application/json",
	}
}
