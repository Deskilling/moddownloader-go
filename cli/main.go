package cli

import (
	"fmt"

	"github.com/deskilling/moddownloader-go/util"
)

func CliMain() {
	fmt.Println("🚀 Welcome to Mod Downloader! Choose an option:")
	fmt.Println("[1] 📦 Mod Files")
	fmt.Println("[2] 🎮 Modpack")
	fmt.Println("[3] Create Default Config")

	var option int
	_, err := fmt.Scanln(&option)
	if err != nil {
		panic(err)
	}

	if option == 1 {
		modMain()
	} else if option == 2 {
		modpackMain()
	} else if option == 3 {
		util.CreateConfig()
	}
}
