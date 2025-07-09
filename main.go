package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/deskilling/moddownloader-go/cli"
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

	if len(os.Args) < 2 {
		cli.CliMain()
	} else {
		util.CheckArgs()
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
