package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"slices"
	"sync"
)

func getDownload(extractedInformation []ModVersionInformation, version string, loader string) (string, string, bool, error) {
	// _ is here the index of the current value and i the index of the current position
	for _, v := range extractedInformation {
		// in the slice from i is the GameVersion and Loder if not loop
		if slices.Contains(v.GameVersions, version) && slices.Contains(v.SupportedLoaders, loader) {
			// Runs when there are no files
			if len(v.Files) == 0 {
				return "", "", true, fmt.Errorf("no files available")
			}

			downloadUrl := v.Files[0].URL
			filename := v.Files[0].Filename

			return downloadUrl, filename, true, nil
		}
	}

	return "", "", false, fmt.Errorf("idfk")
}

func downloadMod(modName string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["modVersionInformation"], modName)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch mod version info: %w", err)
	}

	extractedInformation, err := extractVersionInformation(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse mod version info: %w", err)
	}

	downloadUrl, filename, _, err := getDownload(extractedInformation, version, loader)
	if err != nil {
		return "", false, fmt.Errorf("failed to get download for file: %s", filename)
	}

	err = downloadFile(downloadUrl, filepath+filename)
	if err != nil {
		return "", false, fmt.Errorf("failed to download file: %w", err)
	}

	// To convert the Id into the Title
	projectName, err := projectIdToTitle(modName)
	if err != nil {
		return "", false, fmt.Errorf("failed to get project title: %w", err)
	}

	return projectName, true, nil
}

func downloadViaHash(hash string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(modrinthEndpoint["versionFileHash"], hash)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch version info via hash: %w", err)
	}

	extractedInformation, err := extractVersionHashInformation(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse version info via hash: %w", err)
	}

	modName, status, err := downloadMod(extractedInformation.ProjectId, version, loader, filepath)
	if err != nil || !status {
		return "", false, err
	}

	return modName, true, nil
}

func updateAllViaArgs(version string, loader string, outputPath string, sha1Hashes []string, sha512Hashes []string, allFiles []os.DirEntry) {
	fmt.Printf("\n🎮Version: %s\n", version)
	fmt.Printf("🔧Loader: %s\n", loader)

	// To wait for goroutine
	var wg sync.WaitGroup
	// a lock kinda
	var mu sync.Mutex

	fmt.Println("\n📡 Downloading mods...")

	var downloadedMods []string
	var failedMods []string

	bar := progressbar.NewOptions(int(len(sha1Hashes)),
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
	for indexSha1, atIndexSha1 := range sha1Hashes {
		// Increment WaitGroup counter
		wg.Add(1)

		go func(index int, sha1 string) {
			// Decrement counter when goroutine completes
			defer wg.Done()

			modName, status, err := downloadViaHash(sha1, version, loader, outputPath)
			if err != nil || !status {
				modName, status, err = downloadViaHash(sha512Hashes[index], version, loader, outputPath)
				if err != nil || !status {
					mu.Lock()
					if modName == "" {
						modName = string(allFiles[index].Name())
					}
					failedMods = append(failedMods, modName)
					//fmt.Printf("❌ Failed: %s\n", modName)
					mu.Unlock()
					bar.Add(1)
					// Return is used to exit the goroutine
					return
				}
			}
			mu.Lock()
			//fmt.Println("✅ Downloaded:", modName)
			downloadedMods = append(downloadedMods, modName)
			mu.Unlock()
			bar.Add(1)

		}(indexSha1, atIndexSha1)
	}

	// Wait for all downloads to finish
	wg.Wait()

	fmt.Print("\n\n")
	slices.Sort(downloadedMods)
	if len(downloadedMods) != 0 {
		for i := range len(downloadedMods) {
			fmt.Println("✅ Downloaded:", downloadedMods[i])
		}
	}

	slices.Sort(failedMods)
	if len(failedMods) != 0 {
		for i := range len(failedMods) {
			fmt.Printf("❌ Failed: %s\n", failedMods[i])
		}
	}

	fmt.Println("\n\n✅ All downloads completed.")
}
