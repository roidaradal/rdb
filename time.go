package rdb

import (
	"time"
	_ "time/tzdata"
)

const (
	currentTimezone string = "Asia/Singapore"
	standardFormat  string = "2006-01-02 15:04:05"
	timestampFormat string = "060102150405"
)

func TimeNow() time.Time {
	timezone, _ := time.LoadLocation(currentTimezone)
	return time.Now().In(timezone)
}

func DateTimeNow() DateTime {
	return DateTime(TimeNow().Format(standardFormat))
}

func TimestampNow() string {
	return TimeNow().Format(timestampFormat)
}
