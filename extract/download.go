package extract

import (
	"fmt"
	"slices"
)

func GetDownload(extractedInformation []ModVersionInformation, version string, loader string) (string, string, bool, error) {
	for _, v := range extractedInformation {
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
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
