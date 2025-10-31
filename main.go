package main

import (
	"moddownloader/filesystem"
	"moddownloader/modpack"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func Init() {
	util.InitLogger(-4)
	util.ReadConfig()

}

func main() {
	Init()
	c, _ := filesystem.ReadFile("./modrinth.index.json")
	log.Debug(modpack.GetIdsMrpack(modpack.ParseMrpackJson(c)))
}
