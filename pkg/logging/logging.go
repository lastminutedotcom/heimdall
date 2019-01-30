package log

import (
	"fmt"
	"git01.bravofly.com/golang/appfw.git/pkg/logging"
	"os"
)

func Init() {
	AppLog := logging.NewAppLog(os.Stdout)
	AppLog.UseFileOutput("heimdall")
}

func Info(message string, args ...interface{}) {
	logging.Info(fmt.Sprintf(message, args...), nil)
}

func Debug(message string, args ...interface{}) {
	logging.Debug(fmt.Sprintf(message, args...), nil)
}

func Warn(message string, args ...interface{}) {
	logging.Warn(fmt.Sprintf(message, args...), nil)
}

func Error(message string, args ...interface{}) {
	logging.Error(fmt.Sprintf(message, args...), nil)
}

func Fatal(message string, args ...interface{}) {
	logging.Fatal(fmt.Sprintf(message, args...), nil)
}
