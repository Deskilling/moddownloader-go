package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func modpackMain() {
	var outputPath string = "output/"

	err := checkOutputPath(outputPath)
	if err != nil {
		fmt.Printf("❌ Error checking/creating %s:%s\n", outputPath, err)
		return
	}

	var inputPath string = "modpacks/"

	status, err := doesPathExist(inputPath)
	if err != nil {
		fmt.Printf("❌ Error checking/creating %s: %s\n", inputPath, err)
		return
	}

	if status {
		fmt.Println("📂 Folder `modpacks/` exists!")
	} else {
		fmt.Println("📂 Created `modpacks/` folder")
	}

	scanner := bufio.NewScanner(os.Stdin)

	// gets lates minecraft version
	modrinthVersions, err := getReleaseVersions()
	if err != nil {
		return
	}
	latestVersion := modrinthVersions[0].Version

	fmt.Printf("\n🎮 Enter Minecraft version (default: %s) ➔  ", latestVersion)
	scanner.Scan()
	version := scanner.Text()
	if version == "" {

		version = latestVersion
	}

	fmt.Print("🔧 Enter Loader (default: From Modpack) ➔  ")
	scanner.Scan()
	loader := scanner.Text()
	if loader == "" {
		// Bruh
		loader = ""
	}

	fmt.Println("\n📦 Select Modpack:")
	directory, err := getAllFilesFromDirectory(inputPath, ".mrpack")
	if err != nil {
		return
	}

	if len(directory) == 0 {
		fmt.Println("❌ No modpacks found in the directory. Please add .mrpack files to the modpacks folder.")
		return

	} else if len(directory) == 1 {
		fmt.Println("✅ Found one modpack. Using it automatically: " + directory[0].Name())
		inputPath = inputPath + directory[0].Name()

	} else if len(directory) > 1 {
		fmt.Println("🔍 Multiple modpacks found. Please select one:")
		for i, file := range directory {
			fileName := file.Name()
			noExtensionFilename := fileName[:len(fileName)-len(filepath.Ext(fileName))]
			fmt.Printf("[%d] %s\n", i+1, noExtensionFilename)
		}

		scanner.Scan()
		option, _ := strconv.Atoi(scanner.Text())
		inputPath = inputPath + directory[option-1].Name()
	}

	fmt.Println("📂 Extracting modpack...")
	err = extractZip(inputPath, "temp/")
	if err != nil {
		fmt.Println("❌ Error extracting zip:", err)
		return
	}

	modpackContent := readFile("temp/modrinth.index.json")
	err = checkOutputPath(outputPath)
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

	outputFile := outputPath + version + "_" + parsedModpack.Name + ".mrpack"
	os.Create(outputFile)
	err = zipSource("temp/", outputFile)
	if err != nil {
		fmt.Println("❌ Error zipping:", err)
		return
	}

	fmt.Println("✅ Modpack successfully created at: " + outputFile)
	fmt.Println("\n[Enter to exit]")
	scanner.Scan()
}
