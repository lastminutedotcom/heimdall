package cmd

import (
	"github.com/lastminutedotcom/heimdall/pkg/kubernetes"
	"github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/lastminutedotcom/heimdall/pkg/scheduler"
	"os"
)

func Run() {
	log.Init()

	config := readConfig(os.Getenv("CONFIG_PATH"))
	if config.KubeConfig != nil {
		p := config.KubeConfig.MgmtPort
		kubernetes.ConfigureDeployment(p)
	}

	scheduler.Scheduler{
		Config: config,
	}.Start(Orchestrate())
}

func readConfig(filePath string) *model.Config {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("error opening configuration file: %v", err)
		return nil
	}
	defer file.Close()
	config, err := model.ParseConfig(file)
	if err != nil {
		log.Fatal("error parsing configuration JSON: %v", err)
		return nil
	}
	return config
}
