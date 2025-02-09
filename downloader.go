package main

import (
	"fmt"
	"slices"
)

func downloadMod(modName string, version string, loader string, filepath string) (bool, error) {
	err := checkOutputPath(filepath)
	if err != nil {
		return false, err
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

	projectName, err := projectIdToTitle(modName)
	if err != nil {
		return false, err
	}

	fmt.Printf("Version %s not found for %s, Loader: %s", version, projectName, loader)

	return false, nil
}

func downloadViaHash(hash string, version string, loader string, filepath string) (bool, error) {
	checkOutputPath(filepath)

	url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], hash)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return false, err
	}

	extractedInformation, err := extractVersionHashInformation(response)
	if err != nil {
		return false, err
	}

	status, err := downloadMod(extractedInformation.ProjectId, version, loader, filepath)
	if err != nil || !status {
		return false, err
	}

	return status, nil
}
