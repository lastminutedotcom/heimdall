package waf

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/cmd/client"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpWafs struct {
}

func (h HttpWafs) GetWafTriggersBy(zoneID string, since, until time.Time) ([]model.WafTrigger, error) {
	httpResponse, wafResponse, err := h.callWafTrigger(fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/firewall/events?per_page=50", zoneID))
	if err != nil {
		return nil, fmt.Errorf("get WAF: %v", err)
	}

	triggers := make([]model.WafTrigger, 0)
	if httpResponse.StatusCode == http.StatusOK {
		triggers = h.nextWafTriggersBy(wafResponse.WafTriggers, triggers, zoneID, wafResponse.ResultInfo.NextPageId, since, until)
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

func (h HttpWafs) nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time) []model.WafTrigger {
	for _, wafTrigger := range triggers {
		if after(wafTrigger.OccurredAt, since) {
			continue
		}

		if before(wafTrigger.OccurredAt, until) {
			return result
		}

		if in(wafTrigger.OccurredAt, since, until) {
			result = append(result, wafTrigger)
		}
	}

	if nextPageId != "" {
		nextWafTriggers, actualNextPageId, _ := h.getWafTrigger(zoneID, nextPageId)
		return h.nextWafTriggersBy(nextWafTriggers, result, zoneID, actualNextPageId, since, until)
	}

	return result
}
