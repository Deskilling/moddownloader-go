package util

import "github.com/charmbracelet/log"

func InitLogger(level log.Level) {
	log.SetLevel(level)
	log.SetReportCaller(true)
	log.SetPrefix("moddownloader")
}
