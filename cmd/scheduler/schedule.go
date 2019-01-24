package scheduler

import (
	"fmt"
	"git01.bravofly.com/n7/heimdall/cmd/model"
	"gopkg.in/robfig/cron.v2"
	"log"
	"os"
	"os/signal"
)

var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

type Scheduler struct {
	Config *model.Config
}

func (s Scheduler) Start(function func()) {
	cronExpression := fmt.Sprintf("*/%s * * * *", s.Config.CollectEveryMinutes)
	logger.Printf("start collecting data at every %sth minute of the last %s minute ", s.Config.CollectEveryMinutes, s.Config.CollectEveryMinutes)
	c := cron.New()
	c.AddFunc(cronExpression, function)
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	sign := <-sig
	c.Stop()
	fmt.Println("Got signal:", sign)
}
