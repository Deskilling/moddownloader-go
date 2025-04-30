package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkConnection() error {
	_, err := modrinthWebRequest(modrinthEndpoint["default"])
	if err != nil {
		fmt.Println("❌ An error occurred: Please check your internet connection, or maybe the modrinth api is down")
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
		cliMain()
	} else {
		// Maybe Move
		version, loader, input, output, mode := checkArgs()
		if mode == "mods" {
			fmt.Println("📁 Checking input path...")
			input = checkStringValidPath(input)
			fmt.Println("🔍 Calculating hashes for your mods...")
			sha1Hashes, sha512Hashes, allFiles, _ := calculateAllHashesFromDirectory(input)
			fmt.Println("📁 Checking output path...")
			output = checkStringValidPath(output)
			updateAllViaArgs(version, loader, output, sha1Hashes, sha512Hashes, allFiles)

		} else if mode == "modpack" {
			//output = checkStringValidPath(output)
			err := checkOutputPath(output)
			if err != nil {
				fmt.Println("❌ Failed to check output path:", err)
				return
			}

			if loader != "fabric" {
				fmt.Println("😢 Sowy! Only Fabric loader is supported right now >:(")
				return
			}
			inputPath, err := checkMrpack(input)
			if err != nil {
				fmt.Println("❌ Invalid Modpack: File not found or incorrect format")
				return
			}

			fmt.Println("📂 Extracting modpack...")
			err = extractZip(inputPath, "temp/")
			if err != nil {
				fmt.Println("❌ Error extracting zip:", err)
				return
			}

			modpackContent := readFile("temp/modrinth.index.json")
			err = checkOutputPath(output)
			if err != nil {
				fmt.Println("❌ Error checking/creating output folder:", err)
				return
			}

			fmt.Println("🔍 Parsing modpack...")
			parsedModpack, formatedModpack, err := parseModpack(modpackContent, version, loader)
			if err != nil {
				fmt.Println("❌ Error parsing modpack:", err)
				return
			}
			writeFile("temp/modrinth.index.json", formatedModpack)

			os.Create(output + version + "_" + parsedModpack.Name + ".mrpack") //nolint:errcheck
			err = zipSource("temp/", output+parsedModpack.Name+version+".mrpack ")
			if err != nil {
				fmt.Println("❌ Error zipping:", err)
				return
			}
			fmt.Println("✅ Modpack successfully created at: " + output + parsedModpack.Name + version + ".mrpack")
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
