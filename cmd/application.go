package cmd

import (
	"github.com/lastminutedotcom/heimdall/pkg/client"
	"github.com/lastminutedotcom/heimdall/pkg/client/zone"
	"github.com/lastminutedotcom/heimdall/pkg/data_collector"
	"github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/metric"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/marpaia/graphite-golang"
	"strconv"
	"time"
)

func Orchestrate() func(config *model.Config) {
	return func(config *model.Config) {
		graphite, err := graphite.NewGraphite(config.GraphiteConfig.Host, config.GraphiteConfig.Port)
		if err != nil {
			log.Error("error creating Graphite connection: %v", err)
			return
		}

		aggregates, err := data_collector.GetZones(zone.HttpZones{})
		if err != nil {
			log.Error("error getting zone ids: %v", err)
			return
		}

		endDate := time.Now().In(time.UTC)
		everyMinutes, _ := strconv.Atoi(config.CollectEveryMinutes)
		startDate := endDate.Add(time.Duration(everyMinutes) * time.Minute * -1)

		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startDate.Hour(), startDate.Minute()-1, 59, 0, startDate.Location())
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), endDate.Hour(), endDate.Minute(), 0, 0, endDate.Location())

		log.Info("time range from %s to %s", startDate.Format("2006-01-02T15:04:05Z"), endDate.Format("2006-01-02T15:04:05Z"))

		for _, aggregate := range aggregates {
			response := client.CallGraphQlApi(aggregate.ZoneID, startDate, endDate)
			data_collector.GetColocationTotals(aggregate, response)
			data_collector.GetRatelimitTotals(aggregate, response)
			data_collector.GetWafTotals(aggregate, response)
		}

		if err := adaptAndSend(aggregates, graphite); err != nil {
			log.Error("error converting metrics and sending to Graphite: %v", err)
			return
		}
	}
}

func adaptAndSend(aggregates []*model.Aggregate, g *graphite.Graphite) error {
	for _, a := range aggregates {
		metrics := metric.AdaptMetric(a)
		if err := metric.Push(metrics, g); err != nil {
			return err
		}
	}
	return nil
}
