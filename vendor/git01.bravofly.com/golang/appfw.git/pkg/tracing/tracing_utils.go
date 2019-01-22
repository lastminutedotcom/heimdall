package tracing

import (
	"net/http"
	"os"

	"git01.bravofly.com/golang/appfw.git/pkg/properties"
	"github.com/google/uuid"
)

// NewTraceHeaders will create a new http.Header to propagate tracing information to downstream
func NewTraceHeaders() http.Header {
	return make(map[string][]string, 0)
}

// TraceHeadersFrom detects if an http.Request has already tracing headers set and returns them
func TraceHeadersFrom(request *http.Request) http.Header {
	if request == nil {
		return NewTraceHeaders()
	}
	tracingHeaders := NewTraceHeaders()
	for header := range tracingHeaders {
		if value, found := request.Header[header]; found {
			tracingHeaders[header] = value
		}
	}
	return tracingHeaders
}

// NewTraceID generates a new traceId using UUIDv4
func NewTraceID() (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}

// HostID gets the host id as known by the kernel
func HostID() (string, error) {
	return os.Hostname()
}

// AppName should return the application name
func AppName() string {
	return properties.AppName()
}
