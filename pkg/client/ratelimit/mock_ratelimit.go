package ratelimit

import (
	"encoding/json"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type MockRateLimitClient struct {
	Path string
}

func (m MockRateLimitClient) GetSecurityEvents(zoneID string, since, until time.Time) ([]model.RateLimit, error) {
	file, _ := os.Open(m.Path)
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	rateLimitResponse := model.RateLimitResponse{}
	json.Unmarshal([]byte(byteValue), &rateLimitResponse)
	return rateLimitResponse.Result, nil

}
func (m MockRateLimitClient) callSecurityEvent(url string) (*http.Response, model.RateLimitResponse, error) {

	return nil, model.RateLimitResponse{}, nil
}

func (m MockRateLimitClient) nextSecurityEventsBy(limits []model.RateLimit, result []model.RateLimit, zoneID string, nextPageId string, since, until time.Time) []model.RateLimit {
	return nil

}
func (m MockRateLimitClient) getSecurityEvent(zoneID string, nextPageId string) ([]model.RateLimit, string, error) {

	return nil, "", nil
}
