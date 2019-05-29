package cmd

import (
	"github.com/lastminutedotcom/heimdall/pkg/client/colocation"
	"github.com/lastminutedotcom/heimdall/pkg/client/ratelimit"
	"github.com/lastminutedotcom/heimdall/pkg/client/waf"
	"github.com/lastminutedotcom/heimdall/pkg/client/zone"
	"github.com/lastminutedotcom/heimdall/pkg/data_collector"
	"github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/metric"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/marpaia/graphite-golang"
)

func Orchestration() func(config *model.Config) {
	return func(config *model.Config) {
		newGraphite, err := graphite.NewGraphite(config.GraphiteConfig.Host, config.GraphiteConfig.Port)
		if err != nil {
			log.Error("error creating graphite connection. %v", err)
		}

		httpZonesClient := zone.HttpZones{}
		colocationsClient := colocation.HttpColocations{Config: config}
		httpWafsClient := waf.HttpWafs{}
		httpRateLimitClient := ratelimit.HttpRateLimitClient{}

		aggregate := dataCollector(config, httpZonesClient, colocationsClient, httpWafsClient, httpRateLimitClient)

		metrics := adaptToMetrics(aggregate)
		metric.PushMetrics(metrics, newGraphite)
	}
}

func adaptToMetrics(aggregate []*model.Aggregate) []graphite.Metric {
	return metric.AdaptDataToMetrics(aggregate)

}

func dataCollector(config *model.Config, zoneClient zone.ZonesClient,
	colocationsClient colocation.ColocationsClient, wafClient waf.WafsClient,
	rateLimitClient ratelimit.RateLimitClient) []*model.Aggregate {

	aggregate, _ := data_collector.GetZones(zoneClient)
	data_collector.GetColocationTotals(aggregate, colocationsClient)
	data_collector.GetWafTotals(aggregate, config, wafClient)
	data_collector.GetRatelimitTotals(aggregate, config, rateLimitClient)
	return aggregate
}
