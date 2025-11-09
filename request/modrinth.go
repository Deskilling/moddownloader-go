package request

import (
	"encoding/json"

	"github.com/charmbracelet/log"
)

type Version struct {
	Version     string `json:"version"`
	VersionType string `json:"version_type"`
	//Major       bool   `json:"major"`

	// Fabric Specific
	Build  int  `json:"build"`
	Stable bool `json:"stable"`
}

func GetReleaseVersions() ([]Version, error) {
	versionsData, err := GetBody(ModrinthEndpoint["availableVersions"])
	if err != nil {
		log.Error("error fetching versions", "err", err)
		return nil, err
	}

	var versions []Version
	err = json.Unmarshal([]byte(versionsData), &versions)
	if err != nil {
		log.Error("error unmarshalling versions data", "err", err)
		return nil, err
	}

	var releaseVersions []Version
	for _, v := range versions {
		if v.VersionType == "release" {
			releaseVersions = append(releaseVersions, v)
		}
	}

	return releaseVersions, nil
}
