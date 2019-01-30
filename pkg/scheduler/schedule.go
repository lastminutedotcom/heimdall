package scheduler

import (
	"fmt"
	"git01.bravofly.com/n7/heimdall/pkg/logging"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"gopkg.in/robfig/cron.v2"
	"os"
	"os/signal"
)

type Scheduler struct {
	Config *model.Config
}

func (s Scheduler) Start(function func(config *model.Config)) {
	cronExpression := fmt.Sprintf("*/%s * * * *", s.Config.CollectEveryMinutes)
	log.Info(fmt.Sprintf("start collecting data at every %sth minute of the last %s minute ", s.Config.CollectEveryMinutes, s.Config.CollectEveryMinutes), nil)
	c := cron.New()
	c.AddFunc(cronExpression, func() { function(s.Config) })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	sign := <-sig
	c.Stop()
	fmt.Println("got signal:", sign)
}
