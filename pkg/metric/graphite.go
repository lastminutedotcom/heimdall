package metric

import (
	"git01.bravofly.com/n7/heimdall/pkg/logging"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"github.com/marpaia/graphite-golang"
)

func PushMetrics(aggregate []*model.Aggregate, config *model.Config) {

	metrics := adaptDataToMetrics(aggregate)

	newGraphite, err := graphite.NewGraphite(config.GraphiteConfig.Host, config.GraphiteConfig.Port)

	if err != nil {
		logging.Error("error creating graphite connection. %v", err)
	}

	newGraphite.SendMetrics(metrics)
}
