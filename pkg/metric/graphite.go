package metric

import (
	"github.com/marpaia/graphite-golang"
)

func Push(metrics []graphite.Metric, toServer *graphite.Graphite) error {
	return toServer.SendMetrics(metrics)
}
