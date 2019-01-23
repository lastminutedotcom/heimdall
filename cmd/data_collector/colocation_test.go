package data_collector

import (
	"encoding/json"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"github.com/cloudflare/cloudflare-go"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func readAnalyticsColocationResponse() model.ZoneAnalyticsColocationResponse {
	file, _ := os.Open(filepath.Join("..", "..", "test", "cloudflare_colocation.json"))

	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	analyticsColocationResponse := model.ZoneAnalyticsColocationResponse{}
	json.Unmarshal([]byte(byteValue), &analyticsColocationResponse)
	return analyticsColocationResponse
}

func Test_colocationDataCollection(t *testing.T) {
	aggregate := model.NewAggregate(cloudflare.Zone{
		ID:   "123",
		Name: "zone",
	})

	collectColocation(readAnalyticsColocationResponse().Result, aggregate)

	assert.Equal(t, len(aggregate.Totals), 5)
	key, _ := time.Parse(time.RFC3339, "2019-01-23T15:01:00Z")
	assert.Equal(t, aggregate.Totals[key].RequestAll.Value, 1033)
	assert.Equal(t, aggregate.Totals[key].BandwidthCached.Value, 2821897)
	assert.Equal(t, aggregate.Totals[key].HTTPStatus["2xx"].Value, 861)

}
