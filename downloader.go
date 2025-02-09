package main

import (
	"fmt"
	"slices"
)

func downloadMod(modName string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["modVersionInformation"], modName)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false,  err
	}

	extractedInformation, err := extractVersionInformation(response)
	if err != nil {
		return "", false, err
	}
	// _ is here the index of the current value and i the index of the current position
	for _, v := range extractedInformation {
		// in the slice from i is the GameVersion and Loder if not loop
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
			downloadUrl := v.Files[0].URL
			filename := v.Files[0].Filename

			downloadFromUrl(downloadUrl, filepath+filename)
			return filename, true, nil
		}
	}

	// To convert the Id into the Title 
	projectName, err := projectIdToTitle(modName)
	if err != nil {
		return "", false ,err
	}

	fmt.Printf("Version %s not found for %s, Loader: %s\n", version, projectName, loader)

	return projectName, false, nil
}

func downloadViaHash(hash string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], hash)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false, err
	}

	extractedInformation, err := extractVersionHashInformation(response)
	if err != nil {
		return "", false, err
	}

	modname, status, err := downloadMod(extractedInformation.ProjectId, version, loader, filepath)
	if err != nil || !status {
		return "", false, err
	}

	return modname, status, nil
}
