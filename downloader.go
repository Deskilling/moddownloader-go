package main

import (
	"fmt"
	"slices"
)

func getDownload(extractedInformation []ModVersionInformation, version string, loader string) (string, string, error) {
	// _ is here the index of the current value and i the index of the current position
	for _, v := range extractedInformation {
		// in the slice from i is the GameVersion and Loder if not loop
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
			// Runs when there are no files
			if len(v.Files) == 0 {
				return "", "", fmt.Errorf("no files available")
			}

			downloadUrl := v.Files[0].URL
			filename := v.Files[0].Filename

			return downloadUrl, filename, nil
		}
	}

	return "", "", fmt.Errorf("idfk")
}

func downloadMod(modName string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["modVersionInformation"], modName)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch mod version info: %w", err)
	}

	extractedInformation, err := extractVersionInformation(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse mod version info: %w", err)
	}

	downloadUrl, filename, err := getDownload(extractedInformation, version, loader)
	if err != nil {
		return "", false, fmt.Errorf("failed to get download for file: %w", filename)
	}

	err = downloadFile(downloadUrl, filepath+filename)
	if err != nil {
		return "", false, fmt.Errorf("failed to download file: %w", err)
	}

	// To convert the Id into the Title
	projectName, err := projectIdToTitle(modName)
	if err != nil {
		return "", false, fmt.Errorf("failed to get project title: %w", err)
	}

	return projectName, false, nil
}

func downloadViaHash(hash string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], hash)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch version info via hash: %w", err)
	}

	extractedInformation, err := extractVersionHashInformation(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse version info via hash: %w", err)
	}

	modName, status, err := downloadMod(extractedInformation.ProjectId, version, loader, filepath)
	if err != nil || !status {
		return "", false, err
	}

	return modName, status, nil
}
