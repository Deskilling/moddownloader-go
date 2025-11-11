package main

import (
	"moddownloader/launcher"
	"moddownloader/util"
)

func Init() {
	util.InitLogger(-4)
	util.ReadConfig()
}

func main() {
	Init()

	mp, err := launcher.ReadModpack()
	if err != nil {
		launcher.CreateModpack(launcher.PrismCurrentVersion(), launcher.PrismCurrentLauncher())
	} else {
		launcher.UpdateModpack(launcher.PrismCurrentVersion(), mp)
	}

	/*
		p := launcher.Mmcpack("./mmc-pack.json")
		launcher.PrimsUpdateVersion("1.10", p)
	*/

	/*
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
	*/
}
