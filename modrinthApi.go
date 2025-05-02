package main

import (
	"encoding/json"
	"fmt"
)

func projectIdToTitle(projectId string) (string, error) {
	url := fmt.Sprintf(modrinthEndpoint["modInformation"], projectId)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch project information: %w", err)
	}

	extractedInformation, err := extractModInformation(response)
	if err != nil {
		return "", fmt.Errorf("failed to parse project information: %w", err)
	}

	return extractedInformation.ProjectTitle, nil
}

type Version struct {
	Version     string `json:"version"`
	VersionType string `json:"version_type"`
	//Major       bool   `json:"major"`

	// Fabric Specific
	Build  int  `json:"build"`
	Stable bool `json:"stable"`
}

func getReleaseVersions() ([]Version, error) {
	versionsData, err := modrinthWebRequest(modrinthEndpoint["availableVersions"])
	if err != nil {
		return nil, fmt.Errorf("error fetching versions: %v", err)
	}

	var versions []Version
	err = json.Unmarshal([]byte(versionsData), &versions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling versions data: %v", err)
	}

	var releaseVersions []Version
	for _, v := range versions {
		if v.VersionType == "release" {
			releaseVersions = append(releaseVersions, v)
		}
	}

	return releaseVersions, nil
}
