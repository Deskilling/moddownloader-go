package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deskilling/moddownloader-go/cli"
	"github.com/deskilling/moddownloader-go/downloader"
	"github.com/deskilling/moddownloader-go/filesystem"
	"github.com/deskilling/moddownloader-go/modpack"
	"github.com/deskilling/moddownloader-go/util"
)

func Init() {
	util.CheckPlatform()

	// TODO - Improve messy main.go
	args, err := util.LoadConfig()
	if err == nil || args != util.GetEmptyConfig() {
		settings, err := util.LoadConfig()
		if err != nil {
			return
		}

		sha1Hashes, sha512Hashes, allFiles, _ := filesystem.CalculateAllHashesFromDirectory(settings.Input)

		switch settings.Mode {
		case "mods":
			downloader.UpdateAllViaArgs(args.Version, args.Loader, args.Output, sha1Hashes, sha512Hashes, allFiles)
		case "modpack":
			jsonData, err := filesystem.CheckMrpack(args.Input)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			modpack.ParseModpack(jsonData, args.Version, args.Loader)
		case "export":
			break
		}
	}
}

func main() {
	Init()

	if len(os.Args) < 2 {
		cli.CliMain()
	} else {
		args := util.CheckArgs()
		sha1Hashes, sha512Hashes, allFiles, _ := filesystem.CalculateAllHashesFromDirectory(args.Input)
		switch args.Mode {
		case "mods":
			downloader.UpdateAllViaArgs(args.Version, args.Loader, args.Output, sha1Hashes, sha512Hashes, allFiles)
		case "modpack":
			jsonData, err := filesystem.CheckMrpack(args.Input)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			modpack.ParseModpack(jsonData, args.Version, args.Loader)
		case "export":
			break
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
