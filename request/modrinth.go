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

type Loader struct {
	Icon                  string   `json:"icon"`
	Name                  string   `json:"name"`
	SupportedProjectTypes []string `json:"supported_project_types"`
}

type EndpointMap map[string]string

var ModrinthEndpoint = EndpointMap{
	"default":               "https://api.modrinth.com",
	"modInformation":        "https://api.modrinth.com/v2/project/%s",
	"modVersionInformation": "https://api.modrinth.com/v2/project/%s/version",
	"versionFileHash":       "https://api.modrinth.com/v2/version_file/%s",
	"versionUpdate":         "https://api.modrinth.com/v2/version_file/{hash}/update",
	"availableVersions":     "https://api.modrinth.com/v2/tag/game_version",
	"availableLoaders":      "https://api.modrinth.com/v2/tag/loader",

	// "search": "https://api.modrinth.com/v2/search",
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

func GetAllLoaders() ([]Loader, error) {
	loaderData, err := GetBody(ModrinthEndpoint["availableLoaders"])
	if err != nil {
		log.Error("error fetching loaders")
		return nil, err
	}

	var loader []Loader
	err = json.Unmarshal([]byte(loaderData), &loader)
	if err != nil {
		log.Error("error umarshaling loader json", "err", err)
		return nil, err
	}

	return loader, nil
}
