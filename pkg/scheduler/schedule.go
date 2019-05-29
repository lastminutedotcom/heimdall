package scheduler

import (
	"fmt"
	"github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"gopkg.in/robfig/cron.v2"
	"os"
	"os/signal"
	"syscall"
)

type Scheduler struct {
	Config *model.Config
}

func (s Scheduler) Start(function func(config *model.Config)) {
	cronExpression := fmt.Sprintf("*/%s * * * *", s.Config.CollectEveryMinutes)
	log.Info("start collecting data at every %sth minute of the last %s minute", s.Config.CollectEveryMinutes, s.Config.CollectEveryMinutes)
	c := cron.New()
	c.AddFunc(cronExpression, func() { function(s.Config) })
	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGABRT)
	sign := <-sig
	c.Stop()
	fmt.Println("got signal:", sign)
}
