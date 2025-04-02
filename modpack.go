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

func parseModpack(jsonData string, version string, loader string) (modpack, []byte) {
    var modpack modpack
    err := json.Unmarshal([]byte(jsonData), &modpack)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return modpack, nil
    }

    fmt.Printf("Modpack Name: %s\n", modpack.Name)
    fmt.Printf("Minecraft Version: %s\n", modpack.Dependencies["minecraft"])
    fmt.Printf("Number of Mods: %d\n", len(modpack.Files))

    modpack.Dependencies["minecraft"] = version
    modpack.Dependencies["fabric-loader"] = getLatestFabricVersion()

    var filesToRemove []int

    for i := 0; i < len(modpack.Files); i++ {
        // Directly Modify
        file := &modpack.Files[i] 
        hashSha1 := file.Hashes.Sha1
        hashSha512 := file.Hashes.Sha512

        url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], hashSha1)
        response, err := modrinthWebRequest(url)
        if err != nil {
            url = fmt.Sprintf(modrinthEndpoint["versionFileHash"], hashSha512)
            response, err = modrinthWebRequest(url)
            if err != nil {
                continue
            }
        }

        extractedHashInformation, err := extractVersionHashInformation(response)
        if err != nil {
            continue
        }

        projectId := extractedHashInformation.ProjectId
        url = fmt.Sprintf(modrinthEndpoint["modVersionInformation"], projectId)
        response, err = modrinthWebRequest(url)
        if err != nil {
            continue
        }

        extractedVersionInformation, err := extractVersionInformation(response)
        if err != nil {
            continue
        }

        downloadUrl, _, _, _ := getDownload(extractedVersionInformation, version, loader)
        if downloadUrl == "" {
            filesToRemove = append(filesToRemove, i)
            continue
        }

        file.DownloadUrl = []string{downloadUrl}

        for _, v := range extractedVersionInformation {
            if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
                if len(v.Files) > 0 {
                    file.Hashes.Sha1 = v.Files[0].Hashes.Sha1
                    file.Hashes.Sha512 = v.Files[0].Hashes.Sha512
                    file.FileSize = v.Files[0].Size
                    file.Path = fmt.Sprintf("mods/%s",v.Files[0].Filename)
                }
            }
        }
    }

    for i := len(filesToRemove) - 1; i >= 0; i-- {
        index := filesToRemove[i]
        // Kuss stackoverflow
        modpack.Files = append(modpack.Files[:index], modpack.Files[index+1:]...)
    }

    jsonOutput, err := json.MarshalIndent(modpack, "", "  ")
    if err != nil {
        fmt.Println("Error marshaling JSON:", err)
        return modpack, nil
    }

    return modpack, jsonOutput
}