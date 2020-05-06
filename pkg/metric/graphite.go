package metric

import (
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/marpaia/graphite-golang"
)

func PushMetrics(metrics []graphite.Metric, graph *graphite.Graphite) {
	err := graph.SendMetrics(metrics)
	if err != nil {
		log.Error("error pushing metrics: %v", err)
	}
}
