package model

import (
	"time"
)

type Response struct {
	Data   Data    `json:"data"`
	Errors []Error `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
}

type Data struct {
	Viewer Viewer `json:"viewer"`
}

type Viewer struct {
	Zones []Zones `json:"zones"`
}

type Zones struct {
	FirewallEventsGroups []FirewallEventsGroup `json:"firewallEventsAdaptiveGroups"`
	HttpRequests1mGroups []HttpRequests1mGroup `json:"httpRequests1mGroups"`
}

type FirewallEventsGroup struct {
	Dimensions FirewallEventDimensions `json:"dimensions"`
}

type HttpRequests1mGroup struct {
	HttpRequestDimensions HttpRequestDimensions `json:"dimensions"`
	HttpRequestSum        HttpRequestSum        `json:"sum"`
}

type FirewallEventDimensions struct {
	Action     string    `json:"action"`
	Host       string    `json:"clientRequestHTTPHost"`
	Method     string    `json:"clientRequestHTTPMethodName"`
	OccurredAt time.Time `json:"datetimeMinute"`
	Source     string    `json:"source"`
}

type HttpRequestDimensions struct {
	DatetimeMinute time.Time `json:"datetimeMinute"`
}

type ResponseStatusMap struct {
	ResponseStatus int `json:"edgeResponseStatus"`
	RequestCount   int `json:"requests"`
}

type HttpRequestSum struct {
	Bytes             int                 `json:"bytes"`
	CachedBytes       int                 `json:"cachedBytes"`
	CachedRequests    int                 `json:"cachedRequests"`
	Requests          int                 `json:"requests"`
	ResponseStatusMap []ResponseStatusMap `json:"responseStatusMap"`
}
