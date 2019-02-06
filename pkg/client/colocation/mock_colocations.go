package colocation

import (
	"encoding/json"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"os"
)

type MockColocations struct {
	Path string
}

func (m MockColocations) GetColosAPI(zoneID string) ([]cloudflare.ZoneAnalyticsColocation, error) {
	file, _ := os.Open(m.Path)
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	analyticsColocationResponse := model.ZoneAnalyticsColocationResponse{}
	json.Unmarshal([]byte(byteValue), &analyticsColocationResponse)
	return analyticsColocationResponse.Result, nil
}
