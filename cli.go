package main

import "fmt"

func cliMain() {
	fmt.Println("🚀 Welcome to Mod Downloader! Choose an option:")
	fmt.Println("[1] 📦 Mod Files")
	fmt.Println("[2] 🎮 Modpack")
	var option int
	_, err := fmt.Scanln(&option)
	if err != nil {
		return
	}

	if option == 1 {
		modMain()

	} else if option == 2 {
		modpackMain()
	}
}
