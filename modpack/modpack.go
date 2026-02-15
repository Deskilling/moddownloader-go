package modpack

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"sync"

	"github.com/deskilling/moddownloader-go/extract"
	"github.com/deskilling/moddownloader-go/request"
)

func ParseModpack(jsonData string, version string, loader string) (Modpack, []byte, error) {
	var modpack Modpack
	err := json.Unmarshal([]byte(jsonData), &modpack)
	if err != nil {
		return modpack, nil, fmt.Errorf("failed to unmarshal modpack: %w", err)
	}

	modpack.Dependencies["minecraft"] = version

	if loader == "" {
		if modpack.Dependencies["fabric-loader"] != "" {
			modpack.Dependencies["fabric-loader"] = request.GetLatestFabricVersion()
			loader = "fabric"
		} else if modpack.Dependencies["forge"] != "" {
			modpack.Dependencies["forge"] = request.GetLatestForgeVersion(version)
			loader = "forge"
		} else if modpack.Dependencies["quilt-loader"] != "" {
			modpack.Dependencies["quilt-loader"] = request.GetLatestQuiltVersion()
			loader = "quilt"
		}
	}

	fmt.Printf("📚 Modpack Name: %s\n", modpack.Name)
	fmt.Printf("🎮 Minecraft Version: %s\n", modpack.Dependencies["minecraft"])
	fmt.Printf("💻 Number of Mods: %d\n", len(modpack.Files))
	fmt.Printf("🔧 Loader: %s\n\n", loader)

	var filesToRemove []int
	var removedMods []string
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < len(modpack.Files); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			file := &modpack.Files[i]
			hashSha1 := file.Hashes.Sha1
			hashSha512 := file.Hashes.Sha512

			url := fmt.Sprintf(request.ModrinthEndpoint["versionFileHash"], hashSha1)
			response, err := request.Request(url)
			if err != nil {
				url = fmt.Sprintf(request.ModrinthEndpoint["versionFileHash"], hashSha512)
				response, err = request.Request(url)
				if err != nil {
					return
				}
			}

			extractedHashInformation, err := extract.VersionHash(response)
			if err != nil {
				return
			}

			projectId := extractedHashInformation.ProjectId
			url = fmt.Sprintf(request.ModrinthEndpoint["modVersionInformation"], projectId)
			response, err = request.Request(url)
			if err != nil {
				return
			}

			extractedVersionInformation, err := extract.Version(response)
			if err != nil {
				return
			}

			downloadUrl, _, _, _ := extract.GetDownload(extractedVersionInformation, version, loader)
			if downloadUrl == "" {
				mu.Lock()
				filesToRemove = append(filesToRemove, i)
				removedMods = append(removedMods, file.Path)
				mu.Unlock()
				return
			}

			file.DownloadUrl = []string{downloadUrl}

			for _, v := range extractedVersionInformation {
				if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
					if len(v.Files) > 0 {
						file.Hashes.Sha1 = v.Files[0].Hashes.Sha1
						file.Hashes.Sha512 = v.Files[0].Hashes.Sha512
						file.FileSize = v.Files[0].Size
						file.Path = fmt.Sprintf("mods/%s", v.Files[0].Filename)
					}
					if file.Path != "" {
						break
					}
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Print("\n\n")

	sort.Sort(sort.Reverse(sort.IntSlice(filesToRemove)))

	for _, index := range filesToRemove {
		modpack.Files = append(modpack.Files[:index], modpack.Files[index+1:]...)
	}

	if len(removedMods) > 0 {
		fmt.Println("❌ Failed: ")
		sort.Strings(removedMods)
		for _, mod := range removedMods {
			fmt.Printf("   - %s\n", mod)
		}
		fmt.Printf("\n📅 Total removed mods: %d\n\n", len(removedMods))
	} else {
		fmt.Println("✅ All mods are compatible with this version!")
	}

	jsonOutput, err := json.MarshalIndent(modpack, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return modpack, nil, fmt.Errorf("failed to marshal modpack into json: %w", err)
	}

	return modpack, jsonOutput, nil
}
