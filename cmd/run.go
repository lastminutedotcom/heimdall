package cmd

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/cmd/client"
	"git01.bravofly.com/n7/heimdall/cmd/data_collector"
	"git01.bravofly.com/n7/heimdall/cmd/kubernetes"
	"git01.bravofly.com/n7/heimdall/cmd/metric"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func Run() {
	kubernetes.Liveness()
	kubernetes.Readiness()

	config := readConfig(os.Getenv("CONFIG_PATH"))
	cronExpression := fmt.Sprintf("*/%s * * * *", config.CollectEveryMinutes)

	logger.Printf("start collecting data at every %sth minute of the last %s minute ", config.CollectEveryMinutes, config.CollectEveryMinutes)

	c := cron.New()
	c.AddFunc(cronExpression, orchestrator(config))

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	s := <-sig
	c.Stop()

	fmt.Println("Got signal:", s)
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

func orchestrator(config *model.Config) func() {
	return func() {
		aggregate := dataCollector(config)
		metric.PushMetrics(aggregate, config)
	}
}

func dataCollector(config *model.Config) []*model.Aggregate {
	aggregate, _ := client.GetZonesId()
	aggregate, _ = data_collector.GetColocationTotals(aggregate, config)
	aggregate, _ = data_collector.GetWafTotals(aggregate, config)

	return aggregate
}
