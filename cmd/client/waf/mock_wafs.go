package waf

import (
	"encoding/json"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type MockWafs struct {
}

func (MockWafs) GetWafTriggersBy(zoneID string, since, until time.Time) ([]model.WafTrigger, error) {
	file, _ := os.Open(filepath.Join("..", "..", "test", "cloudflare_waf.json"))
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

func (MockWafs) nextWafTriggersBy(triggers []model.WafTrigger, result []model.WafTrigger, zoneID, nextPageId string, since, until time.Time) []model.WafTrigger {

	return result
}
