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

func Orchestrate() func(config *model.Config) {
	return func(config *model.Config) {
		graphite, err := graphite.NewGraphite(config.GraphiteConfig.Host, config.GraphiteConfig.Port)
		if err != nil {
			log.Error("error creating Graphite connection: %v", err)
			return
		}

		httpZonesClient := zone.HttpZones{}
		colocationsClient := colocation.HttpColocations{Config: config}
		httpWafsClient := waf.HttpWafs{}
		httpRateLimitClient := ratelimit.HttpRateLimitClient{}

		//TODO use unbuffered channel by sending in chan while collecting
		aggregate := collect(config, httpZonesClient, colocationsClient, httpWafsClient, httpRateLimitClient)
		metricStream := make(chan *model.Aggregate, len(aggregate))
		for _, a := range aggregate {
			metricStream <- a
		}
		close(metricStream)
		if err := adaptAndSend(metricStream, graphite); err!=nil {
			log.Error("error converting metrics and sending to Graphite: %v", err)
			return
		}
	}
}

//TODO work this function to have it stream the data as soon as they are read
// we should switch to the official cloudflare-go client to allow
// concurrent calls to be managed by the client rate-limiter instead of relying on our custom impl
func collect(config *model.Config, zoneClient zone.ZonesClient,
	colocationsClient colocation.ColocationsClient, wafClient waf.WafsClient,
	rateLimitClient ratelimit.RateLimitClient) []*model.Aggregate {

	aggregate, err := data_collector.GetZones(zoneClient)
	if err!=nil {
		log.Error("%v", err)
		return nil
	}
	// we don't throw error here because we logged in the methods
	data_collector.GetColocationTotals(aggregate, colocationsClient)
	data_collector.GetWafTotals(aggregate, config, wafClient)
	data_collector.GetRatelimitTotals(aggregate, config, rateLimitClient)
	return aggregate
}

// better use a buffered channel
func adaptAndSend(aggregates chan *model.Aggregate, g *graphite.Graphite) error {
	for a := range aggregates {
		metrics := metric.AdaptMetric(a)
		if err := metric.Push(metrics, g); err!=nil {
			return err
		}
	}
	return nil
}
