package main

import (
	"encoding/json"
)

// To Extractig the Specific Mod information
type ModInformation struct {
	GameVersions       []string `json:"game_versions"`
	SupportedLoaders   []string `json:"loaders"`
	ProjectId          string   `json:"id"`
	ProjectTitle       string   `json:"title"`
	ProjectDescription string   `json:"description"`
	ProjectUpdated     string   `json:"updated"`
	ProjectDownloads   uint     `json:"downloads"`
	ProjectIconUrl     string   `json:"icon_url"`
}

// To get the Files information
type File struct {
	Hashes struct {
		Sha1   string `json:"sha1"`
		Sha512 string `json:"sha512"`
	} `json:"hashes"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

// To Extract the Specific Version Information
type ModVersionInformation struct {
	GameVersions     []string `json:"game_versions"`
	SupportedLoaders []string `json:"loaders"`
	VersionId        string   `json:"id"`
	ProjectId        string   `json:"project_id"`
	VersionName      string   `json:"name"`
	VersionPublished string   `json:"date_published"`
	ProjectDownloads uint     `json:"downloads"`
	Files            []File   `json:"files"`
}

// To extract all the Relevant Information of a Project from the Endpoint modInformation
func extractModInformation(modData string) (ModInformation, error) {
	var extractedInformation ModInformation
	// Umarshal converts the json data into a Go Struct.
	// byte converts the modData json into a byte slice, this is required for Unmarshal
	// Then there is the pointer to the struct which then searches the byte slice and sets the values in the struct
	err := json.Unmarshal([]byte(modData), &extractedInformation)
	if err != nil {
		return ModInformation{}, err
	}

	return extractedInformation, nil
}

// This can be used for the genral version data and the filehash data
func extractVersionInformation(modVersionData string) ([]ModVersionInformation, error) {
	var extractedModVersionsData []ModVersionInformation
	err := json.Unmarshal([]byte(modVersionData), &extractedModVersionsData)
	if err != nil {
		return nil, err
	}

	return extractedModVersionsData, nil
}

// Sad have to reuse due the slice
func extractVersionHashInformation(modVersionData string) (ModVersionInformation, error) {
	var extractedModVersionsData ModVersionInformation
	err := json.Unmarshal([]byte(modVersionData), &extractedModVersionsData)
	if err != nil {
		return ModVersionInformation{}, err
	}

	return extractedModVersionsData, nil
}
