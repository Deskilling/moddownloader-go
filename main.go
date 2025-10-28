package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/deskilling/moddownloader-go/cli"
	"github.com/deskilling/moddownloader-go/downloader"
	"github.com/deskilling/moddownloader-go/filesystem"
	"github.com/deskilling/moddownloader-go/modpack"
	"github.com/deskilling/moddownloader-go/util"
)

func main() {
	// only macos needs sudo perms ...
	if runtime.GOOS == "darwin" {
		if !util.IsRunningAsRoot() {
			fmt.Println("Trying to Relauch with Sudo")
			util.Relaunch()
			return
		}
	}

	fmt.Println(util.GetTime())

	if len(os.Args) < 2 {
		cli.CliMain()
	} else {
		version, loader, input, output, mode := util.CheckArgs()
		sha1Hashes, sha512Hashes, allFiles, _ := filesystem.CalculateAllHashesFromDirectory(input)
		switch mode {
		case "mods":
			downloader.UpdateAllViaArgs(version, loader, output, sha1Hashes, sha512Hashes, allFiles)
		case "modpack":
			jsonData, err := filesystem.CheckMrpack(input)
			if err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}
			modpack.ParseModpack(jsonData, version, loader)
		case "export":
			break
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
