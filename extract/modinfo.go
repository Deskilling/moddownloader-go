package extract

import (
	"encoding/json"

	"github.com/charmbracelet/log"
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

func Mod(modData string) (*ModInformation, error) {
	var mInfo ModInformation
	if err := json.Unmarshal([]byte(modData), &mInfo); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil, err
	}

	return &mInfo, nil
}

func AllVersionHash(modVersionData string) (*[]ModVersionInformation, error) {
	var vInfo []ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &vInfo); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil, err
	}

	return &vInfo, nil
}

func Version(modVersionData string) (*ModVersionInformation, error) {
	var fInfo ModVersionInformation
	if err := json.Unmarshal([]byte(modVersionData), &fInfo); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil, err
	}

	return &fInfo, nil
}
