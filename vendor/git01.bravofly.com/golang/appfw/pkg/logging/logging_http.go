package logging

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"git01.bravofly.com/golang/appfw/pkg/tracing"
)

// writeAccessLogLine is a helper to log HTTP request and response using the conventional standard
func (l *Log) writeAccessLogLine(r *http.Request, rx *responseCatcher) {
	xff := getContextOrHeaderOrDefault(r.Context(), r.Header, "X-Forwarded-For", "-")
	xbftraceid := getContextOrHeaderOrDefault(r.Context(), r.Header, tracing.TraceHeaderTraceID, "-")
	xbfparent := getContextOrHeaderOrDefault(r.Context(), r.Header, tracing.TraceHeaderParentSpanID, "-")
	xbfspanid := getContextOrHeaderOrDefault(r.Context(), r.Header, tracing.TraceHeaderSpanID, "-")
	l.Write(r.RemoteAddr,
		fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		strconv.Itoa(rx.status),
		strconv.Itoa(rx.bytes),
		strconv.Itoa(int(rx.exectime)), //%D in logback??? execution time?
		r.Referer(),
		r.UserAgent(),
		r.Host,
		xff,
		xbftraceid,
		xbfparent,
		xbfspanid)
}

func getContextOrHeaderOrDefault(ctx context.Context, header http.Header, key string, defaultStr string) string {
	if s, ok := ctx.Value(key).(string); ok {
		return s
	}
	if h, ok := header[key]; ok {
		return h[0]
	}
	return defaultStr
}

// HTTPLoggingHandler is the wrapper around any HTTP request that supports logging information on the provided logger
type HTTPLoggingHandler struct {
	traceLog,
	accessLog *Log
	handler http.Handler
}

// ServeHTTP implements http.Handler interface using the logger to get the response from the application handler
func (h HTTPLoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create tracing
	requestMethod := r.RequestURI
	traceRequestEvent := tracing.NewHTTPEvent(r.Context(), make(map[string][]string, 0), tracing.TraceEventServerReceiving)
	ctx, head, evt := traceRequestEvent.Trace()
	// Log incoming request
	h.traceLog.Trace(ctx, head, evt, requestMethod)
	// Create response interceptor
	start := time.Now()
	responseDetails := &responseCatcher{
		w: w,
	}
	//Serve request with tracing
	ctxRequest := r.WithContext(ctx)
	h.handler.ServeHTTP(responseDetails, ctxRequest)
	// Get response details
	responseDetails.exectime = time.Since(start)
	txRespType := tracing.TraceEventServerSending
	if responseDetails.status > 399 {
		txRespType = tracing.TraceEventServerError
	}
	// Trace and log response
	traceResponseEvent := tracing.NewHTTPEvent(ctxRequest.Context(), responseDetails.Header(), txRespType)
	ctx, head, evt = traceResponseEvent.Trace()
	h.traceLog.Trace(ctx, head, evt, requestMethod)
	// Log response
	h.accessLog.writeAccessLogLine(ctxRequest, responseDetails)
}

// NewHTTPLoggingHandler creates a new handler with the provided access and trae loggers and the next application handler to serve requests
func NewHTTPLoggingHandler(accessLog *Log, traceLog *Log, handler http.Handler) http.Handler {
	return &HTTPLoggingHandler{
		traceLog:  traceLog,
		accessLog: accessLog,
		handler:   handler,
	}
}

// responseCatcher is a ResponseWriter interceptor that is used to provide access logging information from the lower handler to the logging handler
type responseCatcher struct {
	w http.ResponseWriter
	status,
	bytes int
	exectime time.Duration
}

func (r *responseCatcher) Header() http.Header {
	return r.w.Header()
}

func (r *responseCatcher) Write(message []byte) (int, error) {
	bytes, err := r.w.Write(message)
	r.bytes += bytes
	return bytes, err
}

func (r *responseCatcher) WriteHeader(status int) {
	r.w.WriteHeader(status)
	r.status = status
}
