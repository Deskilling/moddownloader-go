package request

import (
	"encoding/json"
	"fmt"

	"github.com/deskilling/moddownloader-go/extract"
)

func GetReleaseVersions() ([]Version, error) {
	versionsData, err := Request(ModrinthEndpoint["availableVersions"])
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

func ProjectIdToTitle(projectId string) (string, error) {
	url := fmt.Sprintf(ModrinthEndpoint["modInformation"], projectId)
	response, err := Request(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch project information: %w", err)
	}

	extractedInformation, err := extract.Mod(response)
	if err != nil {
		return "", fmt.Errorf("failed to parse project information: %w", err)
	}

	return extractedInformation.ProjectTitle, nil
}
