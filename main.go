package main

import (
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func Init() {
	util.InitLogger(-4)
	settings := util.ReadConfig()

	log.Debug(settings)
}

func main() {
	Init()
}
