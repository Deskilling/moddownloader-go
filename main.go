package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deskilling/moddownloader-go/cli"
	"github.com/deskilling/moddownloader-go/util"
)

func main() {
	if len(os.Args) < 2 {
		cli.CliMain()
	} else {
		util.CheckArgs()
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
