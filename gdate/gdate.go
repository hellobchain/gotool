package gdate

import "time"

const (
	StdLayout     = "2006-01-02 15:04:05"
	DateLayout    = "2006-01-02"
	TimeLayout    = "15:04:05"
	ISO8601Layout = "2006-01-02T15:04:05-07:00"
)

func Now() time.Time                              { return time.Now() }
func Today() time.Time                            { return time.Now().Truncate(24 * time.Hour) }
func Format(t time.Time, layout string) string    { return t.Format(layout) }
func Parse(str, layout string) (time.Time, error) { return time.Parse(layout, str) }
func FormatStd(t time.Time) string                { return t.Format(StdLayout) }
func ParseStd(str string) (time.Time, error)      { return time.Parse(StdLayout, str) }
