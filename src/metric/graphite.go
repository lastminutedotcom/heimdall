package metric

import (
	"git01.bravofly.com/n7/heimdall/src/model"
	"github.com/marpaia/graphite-golang"
)

func PushMetrics(aggregate []*model.Aggregate) {

	newGraphite, err := graphite.NewGraphite("10.120.172.134", 2113)

	if err != nil {
		newGraphite = graphite.NewGraphiteNop("10.120.172.134", 2113)
	}

	metrics := adaptDataToMetrics(aggregate)
	newGraphite.SendMetrics(metrics)
}
