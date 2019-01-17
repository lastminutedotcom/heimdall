package data_collector

import (
	"git01.bravofly.com/n7/heimdall/src/client"
	"git01.bravofly.com/n7/heimdall/src/model"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func GetWafTotals(aggregates []*model.Aggregate) ([]*model.Aggregate, error) {
	for _, aggregate := range aggregates {
		waf, err := client.GetWafTriggersBy(aggregate.ZoneID)

	}

}
