package timer

import "time"

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	DefaultDateFormat = "2006-01-02"
)

func New() time.Time {
	return time.Now().UTC()
}

func ParseTime(s string) (time.Time, error) {
	return time.Parse(DefaultTimeFormat, s)
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse(DefaultDateFormat, s)
}

func BeginOfDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.UTC().Location())
}

func EndOfDate(t time.Time) time.Time {
	oneDay := 24 * time.Hour
	return BeginOfDate(t).Add(oneDay)
}
