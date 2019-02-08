package ratelimit

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/N7/heimdall/pkg/client"
	"git01.bravofly.com/N7/heimdall/pkg/dates"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpRateLimitClient struct {
	Config *model.Config
}

func (h HttpRateLimitClient) GetSecurityEvents(zoneID string, since, until time.Time) ([]model.RateLimit, error) {
	url := fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/security/events?limit=1000&source=rateLimit&since=%s&until=%s", zoneID, since.Format(time.RFC3339), until.Format(time.RFC3339))
	httpResponse, response, err := h.callSecurityEvent(url)

	if err != nil {
		return nil, fmt.Errorf("get Rate limit: %v", err)
	}

	rateLimits := make([]model.RateLimit, 0)
	if httpResponse.StatusCode == http.StatusOK {
		rateLimits = h.nextSecurityEventsBy(response.Result, rateLimits, zoneID, response.ResultInfo.Cursors.After, since, until)
		return rateLimits, nil
	}
	return nil, fmt.Errorf("get Rate limit HTTP error %d", httpResponse.StatusCode)

}
func (h HttpRateLimitClient) callSecurityEvent(url string) (*http.Response, model.RateLimitResponse, error) {
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.DoHttpCall(request)
	if err != nil {
		return resp, model.RateLimitResponse{}, fmt.Errorf("get Rate limit HTTP call error: %v", err)
	}

	response := model.RateLimitResponse{}
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &response); err != nil {
		return resp, model.RateLimitResponse{}, fmt.Errorf("HTTP body marshal to JSON error: %v", err)
	}

	return resp, response, nil
}

func (h HttpRateLimitClient) nextSecurityEventsBy(limits []model.RateLimit, result []model.RateLimit, zoneID string, nextPageId string, since, until time.Time) []model.RateLimit {
	for _, limit := range limits {
		if dates.After(limit.OccurredAt, until) {
			continue
		}

		if dates.Before(limit.OccurredAt, since) {
			return result
		}

		if dates.In(limit.OccurredAt, until, since) {
			result = append(result, limit)
		}
	}

	if nextPageId != "" {
		nextRateLimits, actualNextPageId, _ := h.getSecurityEvent(zoneID, nextPageId)
		return h.nextSecurityEventsBy(nextRateLimits, result, zoneID, actualNextPageId, since, until)
	}

	return result

}
func (h HttpRateLimitClient) getSecurityEvent(zoneID string, nextPageId string) ([]model.RateLimit, string, error) {
	httpResponse, response, err := h.callSecurityEvent(fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/security/events?limit=1000&source=rateLimit&cursor=%s", zoneID, nextPageId))
	if err != nil {
		return nil, "", fmt.Errorf("get Rate limit: %v", err)
	}

	if httpResponse.StatusCode == http.StatusOK {
		return response.Result, response.ResultInfo.Cursors.After, nil
	}

	return nil, "", fmt.Errorf("get Rate limit HTTP error %d", httpResponse.StatusCode)
}
