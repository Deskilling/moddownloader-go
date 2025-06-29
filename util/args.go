package util

import (
	"flag"

	"github.com/deskilling/moddownloader-go/request"
)

func CheckArgs() (string, string, string, string, string) {
	var latestVersion, _ = request.GetReleaseVersions()

	// Flags
	argMode := flag.String("mode", "mods", "Select between mods or modpacks")

	// Different default loaders
	var defaultLoader string
	var usage string
	if *argMode == "mods" {
		defaultLoader = "fabric"
		usage = "Loader for Mods"
	} else if *argMode == "modpack" {
		defaultLoader = ""
		usage = "Loader for Modpacks keep empty for automatic detection"
	}
	argLoader := flag.String("loader", defaultLoader, usage)

	argVersion := flag.String("version", latestVersion[0].Version, "Minecraft version")
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

/*
func runArgs() {
	// Maybe Move
	version, loader, input, output, mode := checkArgs()
	if mode == "mods" {
		fmt.Println("📁 Checking input path...")
		input = main.checkStringValidPath(input)
		fmt.Println("🔍 Calculating hashes for your mods...")
		sha1Hashes, sha512Hashes, allFiles, _ := main.calculateAllHashesFromDirectory(input)
		fmt.Println("📁 Checking output path...")
		output = main.checkStringValidPath(output)
		main.updateAllViaArgs(version, loader, output, sha1Hashes, sha512Hashes, allFiles)

	} else if mode == "modpack" {
		//output = checkStringValidPath(output)
		err := main.checkOutputPath(output)
		if err != nil {
			fmt.Println("❌ Failed to check output path:", err)
			return
		}

		if loader != "fabric" && loader != "forge" && loader != "quilt" {
			fmt.Println("😢 Sowy! Only Fabric, Forge and Quilt are supported right now >:(")
			return
		}
		inputPath, err := main.checkMrpack(input)
		if err != nil {
			fmt.Println("❌ Invalid Modpack: File not found or incorrect format")
			return
		}

		fmt.Println("📂 Extracting modpack...")
		// Use a proper temp directory path with platform-specific separator
		tempDir := "temp" + string(filepath.Separator)
		err = main.extractZip(inputPath, tempDir)
		if err != nil {
			fmt.Println("❌ Error extracting zip:", err)
			return
		}

		modpackContent := main.readFile(filepath.Join("temp", "modrinth.index.json"))
		err = main.checkOutputPath(output)
		if err != nil {
			fmt.Println("❌ Error checking/creating output folder:", err)
			return
		}

		fmt.Println("🔍 Parsing modpack...")
		parsedModpack, formatedModpack, err := main.parseModpack(modpackContent, version, loader)
		if err != nil {
			fmt.Println("❌ Error parsing modpack:", err)
			return
		}
		main.writeFile(filepath.Join("temp", "modrinth.index.json"), formatedModpack)

		// Use filepath.Join for cross-platform compatibility
		outputFile := filepath.Join(output, version+"_"+parsedModpack.Name+".mrpack")
		os.Create(outputFile) //nolint:errcheck

		// Use filepath.Join for the source directory
		sourceDir := "temp" + string(filepath.Separator)
		err = main.zipSource(sourceDir, outputFile)
		if err != nil {
			fmt.Println("❌ Error zipping:", err)
			return
		}
		fmt.Printf("✅ Modpack successfully created at: %s\n", outputFile)
	}
}


*/
