package utils

import "time"

func FromRFC3339OrNow(ts string) time.Time {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		now := time.Now().UTC()
		return now
	}
	return t
}
