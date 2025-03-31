package main

import (
	"flag"
	"fmt"
	"os"
)

func checkConnection() error {
	_, err := modrinthWebRequest(modrinthEndpoint["default"])
	if err != nil {
		fmt.Println("An error occured: Please check your internet connection, or maybe the modrinth api is down")
		return err
	}

	return nil
}

func checkArgs() (string, string, string, string, string) {
	var latestVersion, _ = getReleaseVersions()

	// Flags
	argMode := flag.String("mode", "mods", "Select between mods or modpacks")
	argVersion := flag.String("version", latestVersion[0].Version, "Minecraft version")
	argLoader := flag.String("loader", "fabric", "Loader")
	argInputFolder := flag.String("input", "mods_to_update/", "Input file")
	argOutputFolder := flag.String("output", "output/", "Output folder")

	// Parse the command-line flags
	flag.Parse()

	mode := *argMode
	version := *argVersion
	loader := *argLoader
	input := *argInputFolder
	output := *argOutputFolder

	return version, loader, input, output, mode
}

func main() {
	err := checkConnection()
	if err != nil {
		return
	}

	if len(os.Args) < 2 {
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
		}
	} else {
		// Maybe Move
		version, loader, input, output, mode := checkArgs()
		if mode == "mods" {
			input = checkStringValidPath(input)
			sha1Hashes, sha512Hashes, allFiles, _ := calculateAllHashesFromDirectory(input)
			output = checkStringValidPath(output)
			updateAllViaArgs(version, loader, output, sha1Hashes, sha512Hashes, allFiles)
		} else if mode == "modpack" {
			modpackMain()
		}
	}
}
