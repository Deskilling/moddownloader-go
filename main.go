package main

import (
	"moddownloader/filesystem"
	"moddownloader/modpack"
	"moddownloader/util"
)

func Init() {
	util.InitLogger(-4)
	util.ReadConfig()
}

func main() {
	Init()
	if util.GetSettings().Automatic.Toggle {
		for i, v := range util.GetSettings().Automatic.Modpacks {
			if err := modpack.UpdateToml(i, v); err != nil {
				return
			}
		}
	} else {
		modpacks, _ := filesystem.ReadDirectory("./", ".mrpack")

		for _, v := range modpacks {
			modpack.ConvertMrpack(v.Name())
		}
	}
}
