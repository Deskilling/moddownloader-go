package main

import (
	"encoding/json"
	"fmt"
	"slices"
)

type modpack struct {
	Dependencies  map[string]string `json:"dependencies"`
	Files         []file            `json:"files"`
	FormatVersion int               `json:"formatVersion"`
	Game          string            `json:"game"`
	Name          string            `json:"name"`
	VersionId     string            `json:"versionId"`
}

type file struct {
	DownloadUrl []string          `json:"downloads"`
	Env         map[string]string `json:"env"`
	FileSize    int               `json:"fileSize"`
	Hashes      hashes            `json:"hashes"`
	Path        string            `json:"path"`
}

type hashes struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
}

func parseModpack(jsonData string) {
	var version string = "1.19.1"
	var loader string = "fabric"

	var modpack modpack
	err := json.Unmarshal([]byte(jsonData), &modpack)
	if err != nil {
		fmt.Println("Fehler beim Unmarshal:", err)
		return
	}

	fmt.Printf("Modpack Name: %s\n", modpack.Name)
	fmt.Printf("Minecraft Version: %s\n", modpack.Dependencies["minecraft"])
	fmt.Printf("Anzahl der Mods: %d\n", len(modpack.Files))

	for i := range len(modpack.Files) {
		hashSha1 := modpack.Files[i].Hashes.Sha1
		hashSha512 := modpack.Files[i].Hashes.Sha512

		url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], hashSha1)
		response, err := modrinthWebRequest(url)
		if err != nil {
			url = fmt.Sprintf(modrinthEndpoint["versionFileHash"], hashSha512)
			response, err = modrinthWebRequest(url)
			if err != nil {
				continue
			}
		}

		extractedHashInformation, _ := extractVersionHashInformation(response)
		projectId := extractedHashInformation.ProjectId

		url = fmt.Sprintf(modrinthEndpoint["modVersionInformation"], projectId)
		response, _ = modrinthWebRequest(url)
		extractedVersionInformation, _ := extractVersionInformation(response)

		downloadUrl, _, _, _ := getDownload(extractedVersionInformation, version, "fabric")

		var downloadUrlSlice []string
		downloadUrlSlice = append(downloadUrlSlice, downloadUrl)

		modpack.Files[i].DownloadUrl = downloadUrlSlice

		for _, v := range extractedVersionInformation {
			// in the slice from i is the GameVersion and Loder if not loop
			if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
				// Runs when there are no files
				if len(v.Files) == 0 {
					continue
				}

				newHashSha1 := v.Files[0].Hashes.Sha1
				newHashSha512 := v.Files[0].Hashes.Sha512

				fmt.Println(newHashSha1)
				fmt.Println(newHashSha512)
			}
		}
	}
}