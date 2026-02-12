package lib

import (
	"strings"
	"time"
)

type Timezone int

const (
	UTC  Timezone = 0
	WIB  Timezone = 7
	WITA Timezone = 8
	WIT  Timezone = 9
)

// CurrentTime func
// will return in format "2006-01-02 15:04:05"
func CurrentTime(format ...string) string {

	form := TimeFormat()
	if len(format) > 0 {
		form = format[0]
	}
	return time.Now().Format(form)
}

// TimeNow func
func TimeNow() time.Time {
	return time.Now()
}

func Tomorrow() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func NextWeek() time.Time {
	return time.Now().Add(7 * 24 * time.Hour)
}

func TimeZone(zone Timezone, d ...time.Time) time.Time {
	t := time.Now()
	if len(d) > 0 {
		t = d[0]
	}
	return t.Add(time.Duration(zone) * time.Hour)
}

// TimeNow func
func TimeNowPtr(t ...time.Time) *time.Time {
	now := time.Now()
	if len(t) > 0 {
		now = t[0]
	}
	return &now
}

func TimeFormat() string {
	return "2006-01-02 15:04:05"
	// return time.RFC3339
}

func SetTime(date time.Time, hour, minute, second int) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, time.UTC)
}

func IsStarted(t time.Time, d ...time.Time) bool {
	x := TimeNow()
	if len(d) > 0 {
		x = d[0]
	}
	return t.Before(x)
}

func InDate(date, check time.Time) bool {
	return check.After(SetTime(date, 0, 0, 1)) && check.Before(SetTime(date, 23, 59, 59))
}

func AddDate(t time.Time, day int) time.Time {
	return t.AddDate(0, 0, day)
}

func AddMonth(t time.Time, month int) time.Time {
	return t.AddDate(0, month, 0)
}

func AddYear(t time.Time, year int) time.Time {
	return t.AddDate(year, 0, 0)
}

func CompareTime(a, b time.Time) bool {
	if a.IsZero() && b.IsZero() {
		return true
	}
	if a.IsZero() || b.IsZero() {
		return false
	}
	return strings.EqualFold(a.Format(TimeFormat()), b.Format(TimeFormat()))
}
