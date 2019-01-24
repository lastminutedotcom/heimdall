package waf

import (
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"net/http"
	"time"
)

type MockWafs struct {
}

func (MockWafs) GetWafTriggersBy(zoneID string, since, until time.Time) ([]model.WafTrigger, error) {
	return nil, nil
}

func (MockWafs) getWafTrigger(zoneID, nextPageId string) ([]model.WafTrigger, string, error) {
	return nil, "", nil
}

func (MockWafs) callWafTrigger(url string) (*http.Response, model.WAFResponse, error) {

	return nil, model.WAFResponse{}, nil
}

func (MockWafs) nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time) []model.WafTrigger {

	return result
}
