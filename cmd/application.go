package cmd

import (
	"git01.bravofly.com/n7/heimdall/pkg/client/colocation"
	"git01.bravofly.com/n7/heimdall/pkg/client/waf"
	"git01.bravofly.com/n7/heimdall/pkg/client/zone"
	"git01.bravofly.com/n7/heimdall/pkg/data_collector"
	"git01.bravofly.com/n7/heimdall/pkg/metric"
	"git01.bravofly.com/n7/heimdall/pkg/model"
)

func Orchestrator() func(config *model.Config) {
	return func(config *model.Config) {
		aggregate := dataCollector(config)
		metric.PushMetrics(aggregate, config)
	}
}

func dataCollector(config *model.Config) []*model.Aggregate {

	aggregate, _ := data_collector.GetZones(zone.HttpZones{})

	aggregate, _ = data_collector.GetColocationTotals(aggregate, colocation.HttpColocations{
		Config: config,
	})

	aggregate, _ = data_collector.GetWafTotals(aggregate, config, waf.HttpWafs{})

	return aggregate
}
