package main

import (
	"moddownloader/launcher"
	"moddownloader/modrinth"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func Init() {
	util.InitLogger(log.DebugLevel)
	if _, err := util.ReadConfig(); err != nil {
		log.Fatal("failed to read config", "err", err)
	}
}

func main() {
	Init()

	//launcher.DumpDatabase("./app.db")

	var loader string
	var version string

	if launcher.IsPrism() {
		version = launcher.PrismCurrentVersion()
		loader = launcher.PrismCurrentLauncher()
	} else {
		return
	}

	mp, err := launcher.ReadModpack()
	if err != nil {
		launcher.CreateModpack(version, loader)
	} else {
		launcher.UpdateModpack(version, mp)
		modrinth.DownloadAll(mp.Ids, version, loader, mp.Output+"/")
	}

	log.Fatal("mfg")
}
