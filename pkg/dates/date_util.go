package dates

import "time"

func In(occurredAt, since time.Time, until time.Time) bool {
	return (occurredAt.Before(since) || occurredAt.Equal(since)) && occurredAt.After(until)
}

func Before(occurredAt, until time.Time) bool {
	return occurredAt.Equal(until) || occurredAt.Before(until)
}

func After(occuredAt, since time.Time) bool {
	return occuredAt.After(since)
}
