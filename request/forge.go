package request

import (
	"fmt"
	"strings"
)

func GetLatestForgeVersion(version string) string {
	url := fmt.Sprintf("https://files.minecraftforge.net/net/minecraftforge/forge/index_%s.html", version)
	response, err := ModrinthWebRequest(url)
	if err != nil {
		return ""
	}

	content := response

	downloadsIndex := strings.Index(content, `<div class="downloads">`)
	if downloadsIndex == -1 {
		fmt.Println("❌ Could not find downloads section")
		return ""
	}

	smallIndex := strings.Index(content[downloadsIndex:], "<small>")
	if smallIndex == -1 {
		fmt.Println("❌ Could not find version information")
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
