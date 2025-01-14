package rdb

import "time"

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

func DateTimeNowWithExpiry(duration time.Duration) (DateTime, DateTime) {
	now := TimeNow()
	expiry := now.Add(duration)
	return DateTime(now.Format(standardFormat)), DateTime(expiry.Format(standardFormat))
}

func ParseTime(datetime DateTime) (time.Time, error) {
	timezone, _ := time.LoadLocation(currentTimezone)
	return time.ParseInLocation(standardFormat, string(datetime), timezone)
}

func CheckIfExpired(expiry DateTime) bool {
	limit, err := ParseTime(expiry)
	if err != nil {
		return true // default to expired if invalid datetime
	}
	now := TimeNow()
	return now.After(limit)
}
