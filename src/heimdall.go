package main

import "C"
import (
	"git01.bravofly.com/n7/heimdall/src/client"
	"git01.bravofly.com/n7/heimdall/src/data_collector"
	"git01.bravofly.com/n7/heimdall/src/metric"
	"git01.bravofly.com/n7/heimdall/src/model"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func main() {
	//logger.Printf("start collecting data %s", "0 * * * * *")
	//
	//c := cron.New()
	//c.AddFunc("0 * * * * *", orchestrator)
	//
	//go c.Start()
	//sig := make(chan os.Signal)
	//signal.Notify(sig, os.Interrupt, os.Kill)
	//s := <-sig
	//c.Stop()
	//
	//fmt.Println("Got signal:", s)

	orchestrator()
}

func orchestrator() {
	aggregate := dataCollector()
	metric.PushMetrics(aggregate)
}

func dataCollector() []*model.Aggregate {
	aggregate, _ := client.GetZonesId()
	aggregate, _ = data_collector.GetColocationTotals(aggregate)
	//aggregate, _ = data_collector.GetWafTotals(aggregate)

	return aggregate
}
