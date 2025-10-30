package request

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
)

func GetLatestFabricVersion() (string, error) {
	response, err := Request("https://meta.fabricmc.net/v2/versions/loader")
	if err != nil {
		log.Error("failed request", "err", err)
		return "", err
	}

	var fabricVersion []Version
	err = json.Unmarshal([]byte(response), &fabricVersion)
	if err != nil {
		log.Error("failed json umarshal", "err", err)
		return "", err
	}
	return fabricVersion[0].Version, nil
}

func GetLatestQuiltVersion() (string, error) {
	response, err := Request("https://meta.quiltmc.org/v3/versions/loader")
	if err != nil {
		log.Error("failed request", "err", err)
		return "", err
	}

	var quiltVersions []Version
	err = json.Unmarshal([]byte(response), &quiltVersions)
	if err != nil {
		log.Error("failed json umarshal", "err", err)
		return "", err
	}
	return quiltVersions[0].Version, nil
}

func GetLatestForgeVersion(version string) (string, error) {
	url := fmt.Sprintf("https://files.minecraftforge.net/net/minecraftforge/forge/index_%s.html", version)
	response, err := Request(url)
	if err != nil {
		log.Error("failed request", "err", err)
		return "", err
	}

	content := response

	downloadsIndex := strings.Index(content, `<div class="downloads">`)
	if downloadsIndex == -1 {
		log.Error("could not find downloads section")
		return "", err
	}

	smallIndex := strings.Index(content[downloadsIndex:], "<small>")
	if smallIndex == -1 {
		log.Error("could not find version information")
		return "", err
	}

	versionStart := downloadsIndex + smallIndex + 7
	versionEnd := strings.Index(content[versionStart:], "</small>")
	if versionEnd == -1 {
		return "", err
	}

	versionString := content[versionStart : versionStart+versionEnd]

	parts := strings.Split(versionString, " - ")
	if len(parts) != 2 {
		return "", err
	}

	return parts[1], nil
}
