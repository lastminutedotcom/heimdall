/*
// Tracing
//
// In order to be able to follow the request/response path along all the microservices
// we use a tracing implementation with the following rules:
//
// * Client generating a new request creates a new traceId, a baseSpanId as parentSpanId (0000000) and a new spanId
// * Client receiving a response uses traceId, parentSpanId and spanId received
// * Client generating a request as a result of a server request keeps traceId and puts server spanId as parentSpanId, creates new spanId
// * Server receiving a request without traceId generates a new traceId, a baseSpanId as parentSpanId (0000000) and a new spanId
// * Server sending a response uses traceId, parentSpanId and spanId generated in ServerRequest event
// * Client/Server not receiving/sending a response trace an Error event
*/
package tracing

import (
	"context"
	"net/http"
)

const (
	TraceEventClientSending   = "CS"
	TraceEventClientReceived  = "CR"
	TraceEventClientError     = "CE"
	TraceEventServerReceiving = "SR"
	TraceEventServerSending   = "SS"
	TraceEventServerError     = "SE"

	TraceHeaderTraceID      = "X-BF-tracing-traceId"
	TraceHeaderSpanID       = "X-BF-tracing-spanId"
	TraceHeaderParentSpanID = "X-BF-tracing-parent-spanId"
	TraceHeaderAppName      = "X-BF-tracing-appName"
	TraceHeaderHostID       = "X-BF-tracing-hostId"
	TraceHeaderMethodName   = "X-BF-tracing-methodName"
	TraceHeaderPairAppName  = "X-BF-tracing-pair-appName"
	TraceHeaderPairHostID   = "X-BF-tracing-pair-hostId"
)

// TraceEvent represents the tracing event happening when a request is handled by a server or a client
// Use of context.Context will be preferred to httpHeader to propagate the tracing details
// TraceEvent also handles the logic on manipulating the incoming/outgoing tracing headers
type TraceEvent interface {
	Trace() (context.Context, http.Header, string)
	spanID() string
	traceID() (string, error)
	parentSpanID() string
	pairAppName() string
	pairHostID() string
}

// HTTPEvent is a tracing event from the Server side
type HTTPEvent struct {
	ctx       context.Context
	headers   http.Header
	eventType string
}

// NewHTTPEvent creates a new TraceEvent from HTTP item
func NewHTTPEvent(ctx context.Context, headers http.Header, eventType string) TraceEvent {
	return &HTTPEvent{
		ctx:       ctx,
		headers:   headers,
		eventType: eventType,
	}
}

// EventType returns the type of event
func (e *HTTPEvent) Trace() (context.Context, http.Header, string) {
	traceID, err := e.traceID()
	if err != nil {
		traceID, _ = NewTraceID()
	}
	e.ctx = context.WithValue(e.ctx, TraceHeaderTraceID, traceID)
	e.ctx = context.WithValue(e.ctx, TraceHeaderParentSpanID, e.parentSpanID())
	e.ctx = context.WithValue(e.ctx, TraceHeaderSpanID, e.spanID())
	e.ctx = context.WithValue(e.ctx, TraceHeaderAppName, AppName())
	hostID, _ := HostID()
	e.ctx = context.WithValue(e.ctx, TraceHeaderHostID, hostID)
	e.ctx = context.WithValue(e.ctx, TraceHeaderPairAppName, e.pairAppName())
	e.ctx = context.WithValue(e.ctx, TraceHeaderPairHostID, e.pairHostID())

	return e.ctx, e.headers, e.eventType
}

// SpanID creates a SpanID for a HTTPEvent
func (e *HTTPEvent) spanID() string {
	if s := valueFromContextOrHeader(e.ctx, e.headers, TraceHeaderSpanID); s != "" {
		return s
	}
	return NewSpanID()
}

// TraceID handles the traceId in request/response headers
func (e *HTTPEvent) traceID() (string, error) {
	if s := valueFromContextOrHeader(e.ctx, e.headers, TraceHeaderTraceID); s != "" {
		return s, nil
	}
	return NewTraceID()
}

func (e *HTTPEvent) parentSpanID() string {
	if s := valueFromContextOrHeader(e.ctx, e.headers, TraceHeaderParentSpanID); s != "" {
		return s
	}
	return BaseSpanID()
}

func (e *HTTPEvent) pairAppName() string {
	return valueFromContextOrHeader(e.ctx, e.headers, TraceHeaderPairAppName)
}

func (e *HTTPEvent) pairHostID() string {
	return valueFromContextOrHeader(e.ctx, e.headers, TraceHeaderPairHostID)
}

func (e *HTTPEvent) Headers() http.Header {
	return e.headers
}

func (e *HTTPEvent) Context() context.Context {
	return e.ctx
}

func valueFromContextOrHeader(ctx context.Context, headers http.Header, key string) string {
	if ctx != nil {
		if pair, ok := ctx.Value(key).(string); ok {
			return pair
		}
	}
	if pair, ok := headers[key]; ok {
		return pair[0]
	}
	return ""
}
