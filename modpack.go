package main

import (
	"encoding/json"
	"fmt"
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
	Downloads []string          `json:"downloads"`
	Env       map[string]string `json:"env"`
	FileSize  int               `json:"fileSize"`
	Hashes    hashes            `json:"hashes"`
	Path      string            `json:"path"`
}

type hashes struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
}

func getDownloadUrl(modName string, version string, loader string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], modName)

	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch mod version info: %w", err)
	}

	extractedFileInformation, err := extractVersionHashInformation(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse mod version info: %w", err)
	}

	projectId := extractedFileInformation.ProjectId

	url = fmt.Sprintf(modrinthEndpoint["modVersionInformation"], projectId)
	response, err = modrinthWebRequest(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch mod version info: %w", err)
	}

	extractedInformation, err := extractVersionInformation(response)

	fmt.Println(getDownload(extractedInformation, version, loader))

	return "", false, nil
}

func parseModpack(jsonData string) {
	var modpack modpack
	err := json.Unmarshal([]byte(jsonData), &modpack)
	if err != nil {
		fmt.Println("Fehler beim Unmarshal:", err)
		return
	}

	fmt.Printf("Modpack Name: %s\n", modpack.Name)
	fmt.Printf("Minecraft Version: %s\n", modpack.Dependencies["minecraft"])
	fmt.Printf("Anzahl der Mods: %d\n", len(modpack.Files))

	for _, file := range modpack.Files {
		fileHash := file.Hashes.Sha1
		getDownloadUrl(fileHash, "1.17.1", "fabric")
	}

}
