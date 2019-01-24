package colocation

import (
	"encoding/json"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"os"
	"path/filepath"
)

type MockColocations struct {
}

func (MockColocations) GetColosAPI(zoneID string) ([]cloudflare.ZoneAnalyticsColocation, error) {
	file, _ := os.Open(filepath.Join("..", "..", "test", "cloudflare_colocation.json"))
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	analyticsColocationResponse := model.ZoneAnalyticsColocationResponse{}
	json.Unmarshal([]byte(byteValue), &analyticsColocationResponse)
	return analyticsColocationResponse.Result, nil
}
