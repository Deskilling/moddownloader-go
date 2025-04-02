package main

import "fmt"

func modpackMain(path string, version string, loader string) {
	fmt.Println("Modpack")

	doesPathExist("modpacks/")
	extractZip("modpacks/EumelcratPack.mrpack", "temp/")
	apored := readFile("temp/modrinth.index.json")
	checkOutputPath("test/")
	_, deropa := parseModpack(apored, version, loader)
	writeFile("test/test.json", deropa)
}
