package request

import (
	"encoding/json"

	"github.com/charmbracelet/log"
)

func GetReleaseVersions() ([]Version, error) {
	versionsData, err := Request(ModrinthEndpoint["availableVersions"])
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
