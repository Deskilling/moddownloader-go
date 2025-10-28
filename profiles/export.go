package profiles

import (
	"fmt"

	"github.com/deskilling/moddownloader-go/filesystem"
)

func exportModpack(modpackPath string, launcherPath string) {
	// extract to temp/ -> copy overrides to launcher dir -> add the pack name between overrides/{modpackname}/mods
	modpackPath, err := filesystem.CheckMrpack(modpackPath)
	if err != nil {
		fmt.Println("❌ Invalid Modpack: File not found or incorrect format")
		return
	}

	err = filesystem.CheckOutputPath("temp/")
	if err != nil {
		return
	}

	err = filesystem.ExtractZip(modpackPath, "temp/")
	if err != nil {
		fmt.Println("Error extracting zip:", err)
		return
	}

	fmt.Println("Copyy ")
	filesystem.CopyDirectory("temp/", launcherPath)
}
