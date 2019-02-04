package zone

import (
	"encoding/json"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"os"
)

type MockZones struct {
	Path string
}

func (m MockZones) GetZonesId() ([]cloudflare.Zone, error) {
	file, _ := os.Open(m.Path)
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	zoneResponse := cloudflare.ZonesResponse{}
	json.Unmarshal([]byte(byteValue), &zoneResponse)
	return zoneResponse.Result, nil

}
