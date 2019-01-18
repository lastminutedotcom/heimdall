package cmd

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/cmd/client"
	"git01.bravofly.com/n7/heimdall/cmd/data_collector"
	"git01.bravofly.com/n7/heimdall/cmd/metric"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func Run(filePath string) {

	config := readConfig(filePath)

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
		logger.Printf("error reading configuration. %v", err)
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
