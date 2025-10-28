package request

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetLatestFabricVersion() string {
	response, err := Request("https://meta.fabricmc.net/v2/versions/loader")
	if err != nil {
		panic(err)
	}

	var fabricVersion []Version
	err = json.Unmarshal([]byte(response), &fabricVersion)
	if err != nil {
		panic(err)
	}
	return fabricVersion[0].Version
}

func GetLatestQuiltVersion() string {
	response, err := Request("https://meta.quiltmc.org/v3/versions/loader")
	if err != nil {
		panic(err)
	}

	var quiltVersions []Version
	err = json.Unmarshal([]byte(response), &quiltVersions)
	if err != nil {
		panic(err)
	}
	return quiltVersions[0].Version
}

func GetLatestForgeVersion(version string) string {
	url := fmt.Sprintf("https://files.minecraftforge.net/net/minecraftforge/forge/index_%s.html", version)
	response, err := Request(url)
	if err != nil {
		return ""
	}

	content := response

	downloadsIndex := strings.Index(content, `<div class="downloads">`)
	if downloadsIndex == -1 {
		fmt.Println("ould not find downloads section")
		return ""
	}

	smallIndex := strings.Index(content[downloadsIndex:], "<small>")
	if smallIndex == -1 {
		fmt.Println("Could not find version information")
		return ""
	}

	versionStart := downloadsIndex + smallIndex + 7
	versionEnd := strings.Index(content[versionStart:], "</small>")
	if versionEnd == -1 {
		return ""
	}

	versionString := content[versionStart : versionStart+versionEnd]

	parts := strings.Split(versionString, " - ")
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}
