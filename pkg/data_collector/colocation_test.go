package data_collector

import (
	"encoding/json"
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_colocationDataCollection(t *testing.T) {
	aggregate := &model.Aggregate{
		ZoneName: "name",
		ZoneID:   "123",
		Totals:   make(map[time.Time]*model.Counters, 0),
	}

	GetColocationTotals(aggregate, mockRequest(filepath.Join("..", "..", "test", "cloudlfare_graphql_response.json")))

	assert.Equal(t, len(aggregate.Totals), 3)
	key, _ := time.Parse(time.RFC3339, "2020-07-26T10:01:00Z")
	assert.Equal(t, aggregate.Totals[key].RequestAll.Value, 65621)
	assert.Equal(t, aggregate.Totals[key].BandwidthCached.Value, 116641475)
	assert.Equal(t, aggregate.Totals[key].HTTPStatus["2xx"].Value, 59551)

}

func mockRequest(path string) *model.Response {
	file, _ := os.Open(path)
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	response := &model.Response{}
	err := json.Unmarshal(byteValue, response)
	if err != nil {
		log.Info("An error occured: %v", err)
	}
	return response
}
