package gf

import (
	"time"
)

func StringTimestamp(timeLayout string, timetype string) int64 {
	timetpl := "2006-01-02 15:04:05"
	if timetype == "date" {
		timetpl = "2006-01-02"
	} else if timetype == "datetime" {
		timetpl = "2006-01-02 15:04"
	}
	times, _ := time.ParseInLocation(timetpl, timeLayout, time.Local)
	timeUnix := times.Unix()
	return timeUnix
}

func TimestampString(timedata interface{}, timetype string) string {
	timetpl := "2006-01-02 15:04:05"
	if timetype == "date" {
		timetpl = "2006-01-02"
	} else if timetype == "datetime" {
		timetpl = "2006-01-02 15:04"
	}
	return time.Unix(timedata.(int64), 0).Format(timetpl)
}

func NowTimestamp() int64 {
	return time.Now().Unix()
}
