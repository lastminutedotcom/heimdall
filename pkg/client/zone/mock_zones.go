package zone

import (
	"encoding/json"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"os"
	"path/filepath"
)

type MockZones struct {
}

func (MockZones) GetZonesId() ([]cloudflare.Zone, error) {
	file, _ := os.Open(filepath.Join("..", "..", "test", "cloudflare_zone.json"))
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	zoneResponse := cloudflare.ZonesResponse{}
	json.Unmarshal([]byte(byteValue), &zoneResponse)
	return zoneResponse.Result, nil

}
