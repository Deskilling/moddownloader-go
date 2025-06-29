package extract

import (
	"fmt"
	"slices"
)

func GetDownload(extractedInformation []ModVersionInformation, version string, loader string) (string, string, bool, error) {
	// _ is here the index of the current value and i the index of the current position
	for _, v := range extractedInformation {
		// in the slice from i is the GameVersion and Loder if not loop
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
			// Runs when there are no files
			if len(v.Files) == 0 {
				return "", "", true, fmt.Errorf("no files available")
			}

			downloadUrl := v.Files[0].URL
			filename := v.Files[0].Filename

			return downloadUrl, filename, true, nil
		}
	}
	return "", "", false, fmt.Errorf("idfk")
}
