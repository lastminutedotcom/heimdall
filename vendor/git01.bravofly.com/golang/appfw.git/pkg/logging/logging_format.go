package logging

import "time"

func timestamp() string {
	return time.Now().Local().Format("2006-01-02T15:04:05.000-0700")
}

type logFormatter interface {
	separator() string
}

type accessLogFormat string

func (a accessLogFormat) separator() string {
	return " | "
}

type appLogFormat string

func (a appLogFormat) separator() string {
	return " | "
}

type traceLogFormat string

func (t traceLogFormat) separator() string {
	return ";"
}
