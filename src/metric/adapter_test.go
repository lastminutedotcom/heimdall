package metric

import (
	"git01.bravofly.com/n7/heimdall/src/model"
	"github.com/cloudflare/cloudflare-go"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func Test_correctAdapting(t *testing.T) {

	data := make([]*model.Aggregate, 0)
	aggregate := model.NewAggregate(cloudflare.Zone{
		ID:   ":: ID ::",
		Name: ":: Name ::",
	})
	//aggregate.TotalBandwidthAll.Value = 5

	now := time.Now()

	aggregate.Totals[now] = model.NewCounters()

	aggregate.Totals[now].BandwidthAll.Value = 5

	data = append(data, aggregate)

	metrics := adaptDataToMetrics(data)
	assert.Equal(t, 10, len(metrics))
	assert.Equal(t, metrics[3].String(), "cloudflare.new.::_name_::.total.bandwidth.all 5 "+now.Format("2006-01-02 15:04:05"))

}
