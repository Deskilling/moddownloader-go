package main

import (
	"moddownloader/modpack"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func Init() {
	util.InitLogger(-4)
	util.ReadConfig()

	log.Debug(util.GetSettings().General.MaxRoutines)

}

func main() {
	Init()
	modpack.RemoveModpack("schwanz")
	modpack.ExtractMrPack("./modpack.mrpack")
}
