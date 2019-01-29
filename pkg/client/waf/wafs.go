package waf

import (
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"net/http"
	"time"
)

type WafsClient interface {
	GetWafTriggersBy(zoneID string, since, until time.Time) ([]model.WafTrigger, error)

	getWafTrigger(zoneID, nextPageId string) ([]model.WafTrigger, string, error)

	nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time) []model.WafTrigger

	callWafTrigger(url string) (*http.Response, model.WAFResponse, error)
}

//
//
//func GetWafTriggersBy(zoneID string, since, until time.Time) ([]model.WafTrigger, error) {
//	httpResponse, wafResponse, err := callWafTrigger(fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/firewall/events?per_page=50", zoneID))
//	if err != nil {
//		return nil, fmt.Errorf("get WAF: %v", err)
//	}
//
//	triggers := make([]model.WafTrigger, 0)
//	if httpResponse.StatusCode == http.StatusOK {
//		triggers = nextWafTriggersBy(wafResponse.WafTriggers, triggers, zoneID, wafResponse.ResultInfo.NextPageId, since, until)
//		return triggers, nil
//	}
//	return nil, fmt.Errorf("get WAF HTTP error %d", httpResponse.StatusCode)
//}
//
//func getWafTrigger(zoneID, nextPageId string) ([]model.WafTrigger, string, error) {
//	httpResponse, wafResponse, err := callWafTrigger(fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/firewall/events?per_page=50&next_page_id=%s", zoneID, nextPageId))
//	if err != nil {
//		return nil, "", fmt.Errorf("get WAF: %v", err)
//	}
//
//	if httpResponse.StatusCode == http.StatusOK {
//		return wafResponse.WafTriggers, wafResponse.ResultInfo.NextPageId, nil
//	}
//
//	return nil, "", fmt.Errorf("get WAF HTTP error %d", httpResponse.StatusCode)
//}
//
//func callWafTrigger(url string) (*http.Response, model.WAFResponse, error) {
//	request, _ := http.NewRequest(http.MethodGet, url, nil)
//
//	resp, err := client.DoHttpCall(request)
//	if err != nil {
//		return resp, model.WAFResponse{}, fmt.Errorf("get WAF HTTP call error: %v", err)
//	}
//
//	response := model.WAFResponse{}
//	b, _ := ioutil.ReadAll(resp.Body)
//	if err := json.Unmarshal(b, &response); err != nil {
//		return resp, model.WAFResponse{}, fmt.Errorf("HTTP body marshal to JSON error: %v", err)
//	}
//
//	return resp, response, nil
//}
//
//func nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time) []model.WafTrigger {
//	for _, wafTrigger := range triggers {
//		if after(wafTrigger.OccurredAt, since) {
//			continue
//		}
//
//		if before(wafTrigger.OccurredAt, until) {
//			return result
//		}
//
//		if in(wafTrigger.OccurredAt, since, until) {
//			result = append(result, wafTrigger)
//		}
//	}
//
//	if nextPageId != "" {
//		nextWafTriggers, actualNextPageId, _ := getWafTrigger(zoneID, nextPageId)
//		return nextWafTriggersBy(nextWafTriggers, result, zoneID, actualNextPageId, since, until)
//	}
//
//	return result
//}
