package utils

import "time"

const timeFormat = "2006-01-02 15:04:05"

func ParseStringToTime(timeStr string) (time.Time, error) {
	return time.ParseInLocation(timeFormat, timeStr, time.Local)
}
