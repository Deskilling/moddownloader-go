package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/deskilling/moddownloader-go/extract"
	"github.com/deskilling/moddownloader-go/filesystem"
	"github.com/deskilling/moddownloader-go/modpack"
	"github.com/deskilling/moddownloader-go/request"
)

// TODO - Add a progress bar for the download

func downloadMod(modName string, version string, loader string, outputPath string) (string, bool, error) {
	url := fmt.Sprintf(request.ModrinthEndpoint["modVersionInformation"], modName)
	response, err := request.Request(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch mod version info: %w", err)
	}

	extractedInformation, err := extract.Version(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse mod version info: %w", err)
	}

	downloadUrl, filename, _, err := extract.GetDownload(extractedInformation, version, loader)
	if err != nil {
		return "", false, fmt.Errorf("failed to get download for file: %s", filename)
	}

	downloadPath := filepath.Join(outputPath, filename)
	err = request.DownloadFile(downloadUrl, downloadPath)
	if err != nil {
		return "", false, fmt.Errorf("failed to download file: %w", err)
	}

	projectName, err := request.ProjectIdToTitle(modName)
	if err != nil {
		return "", false, fmt.Errorf("failed to get project title: %w", err)
	}

	return projectName, true, nil
}

func downloadViaHash(hash string, version string, loader string, filepath string) (string, bool, error) {
	url := fmt.Sprintf(request.ModrinthEndpoint["versionFileHash"], hash)
	response, err := request.Request(url)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch version info via hash: %w", err)
	}

	extractedInformation, err := extract.VersionHash(response)
	if err != nil {
		return "", false, fmt.Errorf("failed to parse version info via hash: %w", err)
	}

	modName, status, err := downloadMod(extractedInformation.ProjectId, version, loader, filepath)
	if err != nil || !status {
		return "", false, err
	}

	return modName, true, nil
}

func UpdateAllViaArgs(version string, loader string, outputPath string, sha1Hashes []string, sha512Hashes []string, allFiles []os.DirEntry) {
	fmt.Printf("\n🎮Version: %s\n", version)
	fmt.Printf("🔧Loader: %s\n", loader)

	var wg sync.WaitGroup
	var mu sync.Mutex

	fmt.Println("\n📡 Downloading mods...")

	var downloadedMods []string
	var failedMods []string

	for indexSha1, atIndexSha1 := range sha1Hashes {
		wg.Add(1)

		go func(index int, sha1 string) {
			defer wg.Done()

			modName, status, err := downloadViaHash(sha1, version, loader, outputPath)
			if err != nil || !status {
				modName, status, err = downloadViaHash(sha512Hashes[index], version, loader, outputPath)
				if err != nil || !status {
					mu.Lock()
					if modName == "" {
						modName = allFiles[index].Name()
					}
					failedMods = append(failedMods, modName)
					mu.Unlock()
					return
				}
			}
			mu.Lock()
			downloadedMods = append(downloadedMods, modName)
			mu.Unlock()

		}(indexSha1, atIndexSha1)
	}

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

func downloadAllModpack(path, version, loader string) {
	err := filesystem.ExtractZip(path, "temp/")
	if err != nil {
		panic(err)
	}

	jsonData := filesystem.ReadFile("temp/modrinth.index.json")
	modpackData, _, _ := modpack.ParseModpack(jsonData, version, loader)

	for i, v := range modpackData.Files {
		fmt.Println(i)
		fmt.Println(v)

	}
}
