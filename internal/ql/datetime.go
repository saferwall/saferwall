package ql

import "time"

const (
	defaultLayout = "2006-01-02T15:04:05"
)

// strToDatetime will convert an RFC3339 datetime to the time.Time go type
// RFC3339 datetime follows the form : "2012-08-21T16:59:22"
func strToDatetime(date string) (time.Time, error) {
	return time.Parse(defaultLayout, date)
}
