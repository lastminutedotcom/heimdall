package metric

import (
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"github.com/marpaia/graphite-golang"
)

func PushMetrics(aggregate []*model.Aggregate) {

	metrics := adaptDataToMetrics(aggregate)

	newGraphite, err := graphite.NewGraphite("10.120.172.134", 2113)

	if err != nil {
		newGraphite = graphite.NewGraphiteNop("10.120.172.134", 2113)
	}

	newGraphite.SendMetrics(metrics)
}
