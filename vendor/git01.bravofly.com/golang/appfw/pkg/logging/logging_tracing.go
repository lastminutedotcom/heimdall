package logging

import (
	"context"
	"net/http"

	"git01.bravofly.com/golang/appfw/pkg/tracing"
)

func (l *Log) writeTraceLogLine(ctx context.Context, headers http.Header, evtType, requestURI string) {
	// traceID, err := event.TraceID()
	// if err != nil {
	// 	Debug(fmt.Sprintf("error generating traceId for current request: %v", err), nil)
	// }
	// hostID, err := tracing.HostID()
	// if err != nil {
	// 	Debug(fmt.Sprintf("error generating hostID: %v", err), nil)
	// }
	// l.Write(traceID,
	// 	event.ParentSpanID(),
	// 	event.SpanID(),
	// 	event.EventType(),
	// 	hostID,
	// 	tracing.AppName(),
	// 	event.PairHostID(),
	// 	event.PairAppName(),
	// 	event.RequestMethod(),
	// )
	// REDO with new impl
	l.Write(getContextOrDefault(ctx, tracing.TraceHeaderTraceID, ""),
		getContextOrDefault(ctx, tracing.TraceHeaderParentSpanID, ""),
		getContextOrDefault(ctx, tracing.TraceHeaderSpanID, ""),
		evtType,
		getContextOrDefault(ctx, tracing.TraceHeaderHostID, ""),
		getContextOrDefault(ctx, tracing.TraceHeaderAppName, ""),
		getContextOrDefault(ctx, tracing.TraceHeaderPairHostID, ""),
		getContextOrDefault(ctx, tracing.TraceHeaderPairAppName, ""),
		requestURI)
}

func getContextOrDefault(ctx context.Context, key string, defaultStr string) string {
	if s, ok := ctx.Value(key).(string); ok {
		return s
	}
	return defaultStr
}

// Trace writes the tracing info to the log
func (l *Log) Trace(ctx context.Context, headers http.Header, evtType, requestURI string) {
	l.writeTraceLogLine(ctx, headers, evtType, requestURI)
}
