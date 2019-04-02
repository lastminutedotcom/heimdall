package waf

import (
	"encoding/json"
	"git01.bravofly.com/N7/heimdall/pkg/logging"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type MockWafs struct {
	Path string
}

func (m MockWafs) GetWafTriggersBy(zoneID string, since, until time.Time) ([]model.WafTrigger, error) {
	file, err := os.Open(m.Path)
	if err != nil {
		log.Info("%v", err)
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	wafResponse := model.WAFResponse{}
	json.Unmarshal([]byte(byteValue), &wafResponse)
	return wafResponse.WafTriggers, nil
}

func (MockWafs) callWafTrigger(url string) (*http.Response, model.WAFResponse, error) {

	return nil, model.WAFResponse{}, nil
}

func (MockWafs) getWafTrigger(zoneID, nextPageId string) ([]model.WafTrigger, string, error) {
	return nil, "", nil
}

func (MockWafs) nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time, callCount int) []model.WafTrigger {

	return result
}
