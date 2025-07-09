package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/deskilling/moddownloader-go/downloader"
	"github.com/deskilling/moddownloader-go/filesystem"
	"github.com/deskilling/moddownloader-go/request"
)

func modMain() {
	var outputPath string = "output/"

	err := filesystem.CheckOutputPath(outputPath)
	if err != nil {
		fmt.Printf("❌ Error checking/creating %s:%s\n", outputPath, err)
		return
	}

	var inputPath string = "mods_to_update/"

	status, err := filesystem.DoesPathExist(inputPath)
	if err != nil {
		fmt.Printf("❌ Error checking/creating %s: %s\n", inputPath, err)
		return
	}

	if status {
		fmt.Println("📂 Folder `mods_to_update/` exists!")
	} else {
		fmt.Println("📂 Created `mods_to_update/`")
	}

	scanner := bufio.NewScanner(os.Stdin)

	// gets lates minecraft version
	modrinthVersions, err := request.GetReleaseVersions()
	if err != nil {
		panic(err)
	}
	latestVersion := modrinthVersions[0].Version

	fmt.Printf("\n🎮 Enter Minecraft version (default: %s) ➔  ", latestVersion)
	scanner.Scan()
	version := scanner.Text()
	if version == "" {
		version = latestVersion
	}

	fmt.Print("🔧 Enter Loader (default: Fabric) ➔  ")
	scanner.Scan()
	loader := scanner.Text()
	if loader == "" {
		loader = "fabric"
	}

	fmt.Printf("\n📥 Please place all mods into ''%s'' and press ENTER↩️  to continue:", inputPath)
	scanner.Scan()

	fmt.Println("🔍 Calculating hashes for your mods...⌛")
	sha1Hashes, sha512Hashes, allFiles, err := filesystem.CalculateAllHashesFromDirectory(inputPath)
	if err != nil {
		fmt.Println("Error calculating file hashes:", err)
		return
	}

	if len(sha1Hashes) != len(sha512Hashes) {
		fmt.Println("⚠️ Hash lists are mismatched! Something went wrong.")
		return
	} else {
		fmt.Printf("✅ Found %d mods to update!\n\n", len(sha1Hashes))
	}

	downloader.UpdateAllViaArgs(version, loader, outputPath, sha1Hashes, sha512Hashes, allFiles)
}
