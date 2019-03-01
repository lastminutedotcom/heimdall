package zone

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/N7/heimdall/pkg/client"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"net/http"
)

type HttpZones struct {
}

func (HttpZones) GetZonesId() ([]cloudflare.Zone, error) {
	url := fmt.Sprintf(client.CloudFlareAPIRoot + "zones?per_page=50")
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.DoHttpCall(request)
	if err != nil {
		return nil, fmt.Errorf("get zones HTTP call error: %v", err)
	}

	response := cloudflare.ZonesResponse{}

	b, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &response); err != nil {
		return nil, fmt.Errorf("HTTP body marshal to JSON error: %v", err)
	}
	if resp.StatusCode == http.StatusOK {
		return response.Result, nil
	}
	return nil, fmt.Errorf("get zones HTTP error %d", resp.StatusCode)
}
