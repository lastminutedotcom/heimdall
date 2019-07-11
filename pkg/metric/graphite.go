package metric

import (
	"github.com/marpaia/graphite-golang"
)

func Push(metrics []graphite.Metric, graph *graphite.Graphite) error {
	return graph.SendMetrics(metrics)
}
