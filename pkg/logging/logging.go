package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

var  appLog = logrus.New()

func Init() {
	appLog.SetOutput(os.Stdout)
}

func Info(message string, args ...interface{}) {
	appLog.Info(fmt.Sprintf(message, args...))
}

func Debug(message string, args ...interface{}) {
	appLog.Debug(fmt.Sprintf(message, args...))
}

func Warn(message string, args ...interface{}) {
	appLog.Warn(fmt.Sprintf(message, args...))
}

func Error(message string, args ...interface{}) {
	appLog.Error(fmt.Sprintf(message, args...))
}

func Fatal(message string, args ...interface{}) {
	appLog.Fatal(fmt.Sprintf(message, args...))
}
