package main

import (
	"encoding/json"
	"fmt"
)

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

type File struct {
	Hashes struct {
		Sha1   string `json:"sha1"`
		Sha512 string `json:"sha512"`
	} `json:"hashes"`
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

func extractModInformation(modData string) (ModInformation, error) {
	var info ModInformation
	// Umarshal converts the json data into a Go Struct.
	// byte converts the modData json into a byte slice, this is required for Unmarshal
	// Then there is the pointer to the struct which then searches the byte slice and sets the values in the struct
	if err := json.Unmarshal([]byte(modData), &info); err != nil {
		return ModInformation{}, fmt.Errorf("failed to unmarshal mod information: %w", err)
	}

	return info, nil
}

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

func extractVersionInformation(modVersionData string) ([]ModVersionInformation, error) {
	var versionsInfo []ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &versionsInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mod version information: %w", err)
	}

	return versionsInfo, nil
}

func extractVersionHashInformation(modVersionData string) (ModVersionInformation, error) {
	var fileInfo ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &fileInfo); err != nil {
		return ModVersionInformation{}, fmt.Errorf("failed to unmarshal mod version information: %w", err)
	}

	return fileInfo, nil
}
