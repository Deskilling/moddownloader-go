package main

import (
	"fmt"
)

func main() {
	fmt.Println("🚀 Started Moddownloader-go")
	fmt.Println("[1] Modfiles or [2] Modpack")

	var option int = 2
	/*
		_, err := fmt.Scanln(&option)
		if err != nil {
			return
		}
	*/

	if option == 1 {
		modMain()

	} else if option == 2 {
		modpackMain()

	} else {
		fmt.Println("Invalid option")
		main()
	}
}
