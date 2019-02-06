package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	// LogTypeAccess is the file name part for access logs
	LogTypeAccess = "access"
	// LogTypeApplication is the file name part for application logs
	LogTypeApplication = "app"
	// LogTypeTracing is the file name part for tracing logs
	LogTypeTracing = "tracing"
)

// Logging singletons
var AccessLog = NewAccessLog(os.Stdout)
var AppLog = NewAppLog(os.Stdout)
var TraceLog = NewTraceLog(os.Stdout)

// Log is the struct for logging HTTP requests, application mesages or tracing messages
// Depending on the formatter set to the Log instance, the composition of messages will be different
type Log struct {
	formatter logFormatter
	logger    *log.Logger
	logType   string
}

// NewAccessLog creates a new logger for incoming HTTP requests directing output to the Writer
// the standard time format is used as flags (2009/01/23 01:23:23.123123 +0000)
func NewAccessLog(out io.Writer) *Log {
	var format accessLogFormat
	return &Log{
		formatter: format,
		logger:    log.New(out, "", 0),
		logType:   LogTypeAccess,
	}
}

// NewAppLog creates a new logger for incoming HTTP requests directing output to the Writer
// the standard time format is used as flags (2009/01/23 01:23:23.123123 +0000)
func NewAppLog(out io.Writer) *Log {
	var format appLogFormat
	return &Log{
		formatter: format,
		logger:    log.New(out, "", 0),
		logType:   LogTypeApplication,
	}
}

// NewTraceLog creates a new logger for tracing information directing output to the Writer
func NewTraceLog(out io.Writer) *Log {
	var format traceLogFormat
	return &Log{
		formatter: format,
		logger:    log.New(out, "", 0),
		logType:   LogTypeTracing,
	}
}

func (l *Log) format(fields ...string) string {
	return timestamp() + l.formatter.separator() + strings.Join(fields, l.formatter.separator())
}

// Write a log line with multiple fields separated by separator
func (l *Log) Write(fields ...string) {
	l.logger.Println(l.format(fields...))
}

// UseFileOutput configure the log file destination, with rotation and filename using the conventional rules
func (l *Log) UseFileOutput(applicationName string) {
	l.multiplexToStdOutAnd(&lumberjack.Logger{
		Filename: fmt.Sprintf("/appfw/logs/%s-%s.log", l.logType, applicationName),
		// Unlimited daily size
		MaxSize:    0,
		MaxBackups: 15,
		MaxAge:     15,
		Compress:   false,
	})
}

func (l *Log) multiplexToStdOutAnd(out io.Writer) {
	l.logger.SetOutput(io.MultiWriter(os.Stdout, out))
}
