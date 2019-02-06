package logging

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

const (
	LeveLInfo  = "INFO"
	LeveLDebug = "DEBUG"
	LeveLWarn  = "WARN"
	LeveLError = "ERROR"
	LeveLFatal = "FATAL"
)

// Info writes the INFO message in the application logs according to the std format
func Info(message string, req *http.Request) {
	var caller string
	if pc, _, _, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	}
	AppLog.writeAppLogLine(LeveLInfo, message, req, caller)
}

// Debug writes the DEBUG message in the application logs according to the std format
func Debug(message string, req *http.Request) {
	var caller string
	if pc, _, _, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	}
	AppLog.writeAppLogLine(LeveLDebug, message, req, caller)
}

// Warn writes the WARN message in the application logs according to the std format
func Warn(message string, req *http.Request) {
	var caller string
	if pc, _, _, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	}
	AppLog.writeAppLogLine(LeveLWarn, message, req, caller)
}

// Error writes the ERROR message in the application logs according to the std format
func Error(message string, req *http.Request) {
	var caller string
	if pc, _, _, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	}
	AppLog.writeAppLogLine(LeveLError, message, req, caller)
}

// Fatal writes the FATAL message in the application logs and ends the running program, use it with caution
func Fatal(message string, req *http.Request) {
	var caller string
	if pc, _, _, ok := runtime.Caller(1); ok {
		caller = fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
	}
	AppLog.writeAppLogLine(LeveLFatal, message, req, caller)
	os.Exit(1)
}

func (l *Log) writeAppLogLine(level, message string, r *http.Request, caller string) {
	var (
		xbftraceid = "-"
		xbfparent  = "-"
		xbfspanid  = "-"
	)
	if r != nil {
		if h, ok := r.Header["X-BF-tracing-traceId"]; ok {
			xbftraceid = h[0]
		}
		if h, ok := r.Header["X-BF-tracing-parent-spanId"]; ok {
			xbfparent = h[0]
		}
		if h, ok := r.Header["X-BF-tracing-spanId"]; ok {
			xbfspanid = h[0]
		}
	}
	l.Write(level,
		xbftraceid,
		xbfparent,
		xbfspanid,
		"-", // BusinessProfile??
		"-", // Another field...
		caller,
		message,
	)
}
