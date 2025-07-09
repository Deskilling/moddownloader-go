package cli

import "fmt"

func CliMain() {
	fmt.Println("🚀 Welcome to Mod Downloader! Choose an option:")
	fmt.Println("[1] 📦 Mod Files")
	fmt.Println("[2] 🎮 Modpack")

	var option int
	_, err := fmt.Scanln(&option)
	if err != nil {
		panic(err)
	}

	if option == 1 {
		modMain()

	} else if option == 2 {
		modpackMain()
	}
}
