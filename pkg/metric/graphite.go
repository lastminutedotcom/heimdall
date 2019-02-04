package metric

import (
	"github.com/marpaia/graphite-golang"
)

func PushMetrics(metrics []graphite.Metric, graph *graphite.Graphite) {
	graph.SendMetrics(metrics)
}
