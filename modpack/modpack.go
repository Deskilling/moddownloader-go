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
		fmt.Println("❌ Error unmarshalling JSON:", err)
		return modpack, nil, fmt.Errorf("failed to unmarshal modpack: %w", err)
	}

	modpack.Dependencies["minecraft"] = version

	if loader == "" {
		if modpack.Dependencies["fabric-loader"] != "" {
<<<<<<< HEAD:modpack/modpack.go
			modpack.Dependencies["fabric-loader"] = request.GetLatestFabricVersion()
			loader = "fabric"
		} else if modpack.Dependencies["forge"] != "" {
			modpack.Dependencies["forge"] = request.GetLatestForgeVersion(version)
			loader = "forge"
			// not checked but hopefully
		} else if modpack.Dependencies["quilt-loader"] != "" {
			modpack.Dependencies["quilt-loader"] = request.GetLatestQuiltVersion()
			loader = "quilt"
=======
			modpack.Dependencies["fabric-loader"] = getLatestFabricVersion()
			loader = "fabric"
		} else if modpack.Dependencies["forge"] != "" {
			modpack.Dependencies["forge"] = getLatestForgeVersion(version)
			loader = "forge"
>>>>>>> main:modpack.go
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

<<<<<<< HEAD:modpack/modpack.go
=======
	bar := progressbar.NewOptions(len(modpack.Files),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

>>>>>>> main:modpack.go
	for i := 0; i < len(modpack.Files); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			// Directly Modify
			file := &modpack.Files[i]
			hashSha1 := file.Hashes.Sha1
			hashSha512 := file.Hashes.Sha512

			url := fmt.Sprintf(request.ModrinthEndpoint["versionFileHash"], hashSha1)
			response, err := request.ModrinthWebRequest(url)
			if err != nil {
				url = fmt.Sprintf(request.ModrinthEndpoint["versionFileHash"], hashSha512)
				response, err = request.ModrinthWebRequest(url)
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
			response, err = request.ModrinthWebRequest(url)
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

	wg.Wait() // Wait for all goroutines to complete
	fmt.Print("\n\n")

	// Sort filesToRemove to ensure consistent removal order
	sort.Sort(sort.Reverse(sort.IntSlice(filesToRemove)))

	for _, index := range filesToRemove {
		modpack.Files = append(modpack.Files[:index], modpack.Files[index+1:]...)
	}

	/*
	   for i := len(filesToRemove) - 1; i >= 1; i-- {
	       index := filesToRemove[i]
	       fmt.Println(i, index)
	       // Kuss stackoverflow
	       modpack.Files = append(modpack.Files[:index], modpack.Files[index+1:]...)
	   }
	*/

	/*
	   jsonOutput, err := json.Marshal(modpack)
	   if err != nil {
	       fmt.Println("Error marshaling JSON:", err)
	       return modpack, nil, fmt.Errorf("failed to marshal modpack into json: %w", err)
	   }
	*/

	// Print removed mods
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
		fmt.Println("❌ Error marshaling JSON:", err)
		return modpack, nil, fmt.Errorf("failed to marshal modpack into json: %w", err)
	}

	return modpack, jsonOutput, nil
}
