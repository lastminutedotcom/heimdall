package colocation

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/pkg/client"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"net/http"
)

type HttpColocations struct {
	Config *model.Config
}

func (h HttpColocations) GetColosAPI(zoneID string) ([]cloudflare.ZoneAnalyticsColocation, error) {
	url := fmt.Sprintf(client.CloudFlareAPIRoot+"zones/%s/analytics/colos?since=-%s&until=-%s&continuous=%s", zoneID, h.Config.CollectEveryMinutes, "1", "false")
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.DoHttpCall(request)
	if err != nil {
		return nil, fmt.Errorf("get colocation analytics HTTP call error: %v", err)
	}
	response := model.ZoneAnalyticsColocationResponse{}
	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, fmt.Errorf("HTTP body marshal to JSON error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		return response.Result, nil
	}
	return nil, fmt.Errorf("get colocation analytics HTTP error %d", resp.StatusCode)
}
