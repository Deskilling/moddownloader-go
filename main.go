package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file := getLauncherProfiles()
	parseLauncherProfiles(file)

	err := checkConnection()
	if err != nil {
		return
	}

	if len(os.Args) < 2 {
		cliMain()
	} else {
		runArgs()
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
