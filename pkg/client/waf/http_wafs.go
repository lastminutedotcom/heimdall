package waf

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/N7/heimdall/pkg/client"
	"git01.bravofly.com/N7/heimdall/pkg/dates"
	"git01.bravofly.com/N7/heimdall/pkg/logging"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpWafs struct {
}

func (h HttpWafs) GetWafTriggersBy(zoneID string, since, until time.Time, callCount int) ([]model.WafTrigger, error) {
	httpResponse, wafResponse, err := h.callWafTrigger(fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/firewall/events?per_page=50", zoneID))
	if err != nil {
		return nil, fmt.Errorf("get WAF: %v", err)
	}

	triggers := make([]model.WafTrigger, 0)
	if httpResponse.StatusCode == http.StatusOK {
		triggers = h.nextWafTriggersBy(wafResponse.WafTriggers, triggers, zoneID, wafResponse.ResultInfo.NextPageId, since, until, callCount)
		return triggers, nil
	}
	return nil, fmt.Errorf("get WAF HTTP error %d", httpResponse.StatusCode)
}

func (h HttpWafs) getWafTrigger(zoneID, nextPageId string) ([]model.WafTrigger, string, error) {
	httpResponse, wafResponse, err := h.callWafTrigger(fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/firewall/events?per_page=50&next_page_id=%s", zoneID, nextPageId))
	if err != nil {
		return nil, "", fmt.Errorf("get WAF: %v", err)
	}

	if httpResponse.StatusCode == http.StatusOK {
		return wafResponse.WafTriggers, wafResponse.ResultInfo.NextPageId, nil
	}

	return nil, "", fmt.Errorf("get WAF HTTP error %d", httpResponse.StatusCode)
}

func (h HttpWafs) callWafTrigger(url string) (*http.Response, model.WAFResponse, error) {
	log.Info("calling url: %s", url)
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.DoHttpCall(request)
	if err != nil {
		return resp, model.WAFResponse{}, fmt.Errorf("get WAF HTTP call error: %v", err)
	}

	response := model.WAFResponse{}
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &response); err != nil {
		return resp, model.WAFResponse{}, fmt.Errorf("HTTP body marshal to JSON error: %v", err)
	}

	return resp, response, nil
}

func (h HttpWafs) nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time, callCount int) []model.WafTrigger {
	for _, wafTrigger := range triggers {
		if dates.After(wafTrigger.OccurredAt, since) {
			continue
		}

		if dates.Before(wafTrigger.OccurredAt, until) {
			return result
		}

		if dates.In(wafTrigger.OccurredAt, since, until) {
			result = append(result, wafTrigger)
		}
	}

	if nextPageId != "" && callCount < 150 {
		nextWafTriggers, actualNextPageId, _ := h.getWafTrigger(zoneID, nextPageId)
		callCount++
		log.Info("waf page: %d", callCount)
		return h.nextWafTriggersBy(nextWafTriggers, result, zoneID, actualNextPageId, since, until, callCount)
	}

	return result
}
