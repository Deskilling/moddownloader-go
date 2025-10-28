package util

import (
	"flag"

	"github.com/deskilling/moddownloader-go/request"
)

type Args struct {
	Mode    string
	Version string
	Loader  string
	Input   string
	Output  string
}

func CheckArgs() Args {
	var arg Args

	var latestVersion, _ = request.GetReleaseVersions()

	argMode := flag.String("mode", "mods", "Select between mods or modpacks")

	var defaultLoader string
	var usage string

	switch *argMode {
	case "mods":
		defaultLoader = "fabric"
		usage = "Loader for Mods"
	case "modpack":
		defaultLoader = ""
		usage = "Loader for Modpacks keep empty for automatic detection"
	}

	argLoader := flag.String("loader", defaultLoader, usage)

	argVersion := flag.String("version", latestVersion[0].Version, "Minecraft version")
	argInputFolder := flag.String("input", "mods_to_update/", "Input file")
	argOutputFolder := flag.String("output", "output/", "Output folder")

	flag.Parse()

	arg.Mode = *argMode
	arg.Version = *argVersion
	arg.Loader = *argLoader
	arg.Input = *argInputFolder
	arg.Output = *argOutputFolder

	return arg
}
