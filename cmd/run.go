package cmd

import (
	"encoding/json"
	"fmt"
	"git01.bravofly.com/n7/heimdall/pkg/kubernetes"
	"git01.bravofly.com/n7/heimdall/pkg/logging"
	"git01.bravofly.com/n7/heimdall/pkg/model"
	"git01.bravofly.com/n7/heimdall/pkg/scheduler"
	"io/ioutil"
	"os"
)

//var _logger = logging.NewAppLog(os.Stdout)
//var logger = log.New(os.Stdout, "[HEIMDALL] ", log.LstdFlags)

func Run() {
	logging.Init()
	kubernetes.Readiness()
	kubernetes.Liveness()

	config := readConfig(os.Getenv("CONFIG_PATH"))

	scheduler.Scheduler{
		Config: config,
	}.Start(Orchestrator())

}

func readConfig(filePath string) *model.Config {
	file, err := os.Open(filePath)
	if err != nil {
		logging.Error(fmt.Sprintf("error reading configuration. %v", err), nil)
		//logger.Printf("error reading configuration. %v", err)
		return model.DefautConfig()
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var config *model.Config
	json.Unmarshal([]byte(byteValue), &config)
	return config
}
