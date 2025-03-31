package main

import "fmt"

func modpackMain() {
	fmt.Println("Modpack")

	doesPathExist("modpacks/")
	extractZip("modpacks/EumelcratPack.mrpack", "temp/")
	apored := readFile("temp/modrinth.index.json")
	parseModpack(apored)
}
