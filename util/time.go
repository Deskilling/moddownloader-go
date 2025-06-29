package util

import "time"

func GetTime() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02T15:04:05-0700")
	return formattedTime
}
