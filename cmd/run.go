package cmd

import (
	"git01.bravofly.com/N7/heimdall/pkg/kubernetes"
	"git01.bravofly.com/N7/heimdall/pkg/logging"
	"git01.bravofly.com/N7/heimdall/pkg/model"
	"git01.bravofly.com/N7/heimdall/pkg/scheduler"
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
	}.Start(Orchestration())
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
