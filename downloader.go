package main

import (
	"fmt"
	"os"
	"slices"
)

func downloadMod(modName string, version string, loader string, filepath string) (bool, error) {
	exists, err := doesPathExist(filepath)
	if err != nil {
		return false, err
	}

	if exists {
		empty, err := isDirEmpty("output")
		if err != nil {
			return false, err
		}

		if !empty {
			os.RemoveAll(filepath)
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	url := fmt.Sprintf(modrinthEndpoint["modVersionInformation"], modName)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return false, err
	}

	extractedInformation, err := extractVersionInformation(response)
	if err != nil {
		return false, err
	}

	for _, i := range extractedInformation {
		if slices.Contains(i.GameVersions, version) && slices.Contains(i.SupportedLoaders, loader) {
			downloadUrl := i.Files[0].URL
			filename := i.Files[0].Filename

			downloadFromUrl(downloadUrl, filepath+filename)
			return true, nil
		}
	}

	return false, nil
}
