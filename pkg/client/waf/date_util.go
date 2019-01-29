package waf

import "time"

func in(occurredAt, since time.Time, until time.Time) bool {
	return (occurredAt.Before(since) || occurredAt.Equal(since)) && occurredAt.After(until)
}

func before(occurredAt, until time.Time) bool {
	return occurredAt.Equal(until) || occurredAt.Before(until)
}

func after(occuredAt, since time.Time) bool {
	return occuredAt.After(since)
}
