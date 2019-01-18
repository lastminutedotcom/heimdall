package main

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/src/client"
	"git01.bravofly.com/n7/heimdall/src/data_collector"
	"git01.bravofly.com/n7/heimdall/src/metric"
	"git01.bravofly.com/n7/heimdall/src/model"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func main() {

	config := readConfig("config.json")

	logger.Printf("start collecting data %s", config.CronExpression)

	c := cron.New()
	c.AddFunc(config.CronExpression, orchestrator)

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig
	c.Stop()

	fmt.Println("Got signal:", s)

	//orchestrator()
}

func readConfig(filePath string) *model.Config {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Printf("could not open config file %s: %v", file, err)
		return model.DefautConfig()
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var config *model.Config
	json.Unmarshal([]byte(byteValue), &config)
	return config
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
