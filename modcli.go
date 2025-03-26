package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func modMain() {
	var outputPath string = "output/"

	err := checkOutputPath(outputPath)
	if err != nil {
		fmt.Println("❌ Error checking/creating output/:", err)
		return
	}

	status, err := doesPathExist("mods_to_update/")
	if err != nil {
		fmt.Println("❌ Error checking/creating mods_to_update/:", err)
		return
	}

	if status {
		fmt.Println("📂 Folder `mods_to_update/` exists!")
	} else {
		fmt.Println("📂 Created `mods_to_update/`")
	}

	scanner := bufio.NewScanner(os.Stdin)

	// gets lates minecraft version
	modrinthVersions, err := getReleaseVersions()
	if err != nil {
		return
	}
	latestVersion := modrinthVersions[0].Version

	fmt.Printf("\n🎮 Enter Minecraft version (default: %s) ➔  ", latestVersion)
	scanner.Scan()
	version := scanner.Text()
	if version == "" {

		version = latestVersion
	}

	fmt.Print("🔧 Enter Loader (default: Fabric) ➔  ")
	scanner.Scan()
	loader := scanner.Text()
	if loader == "" {
		loader = "fabric"
	}

	fmt.Println("\n" + `📥 Please place all mods into "mods_to_update/" and press ENTER↩️  to continue:`)
	scanner.Scan()

	fmt.Println("🔍 Calculating hashes for your mods...⌛")
	sha1Hashes, sha512Hashes, allFiles, err := calculateAllHashesFromDirectory("mods_to_update/")
	if err != nil {
		fmt.Println("Error calculating file hashes:", err)
		return
	}

	if len(sha1Hashes) != len(sha512Hashes) {
		fmt.Println("⚠️ Hash lists are mismatched! Something went wrong.")
		return
	} else {
		fmt.Printf("✅ Found %d mods to update!\n\n", len(sha1Hashes))
	}

	// Shitty
	// i is the index and v the value at that index
	/*
		for indexSha1, atIndexSha1 := range sha1Hashes {
			modName, downloaded, err := downloadViaHash(atIndexSha1, version, loader, "output/")
			if err != nil || !downloaded {
				modName, downloaded, err := downloadViaHash(sha512Hashes[indexSha1], version, loader, "output/")
				if err != nil || !downloaded {
					fmt.Println("Failed to download")
				} else {
					fmt.Println("Downloaded: ", modName)
				}
			} else {
				fmt.Println("Downloaded: ", modName)
			}
		}
	*/

	// To wait for goroutine
	var wg sync.WaitGroup
	// a lock kinda
	var mu sync.Mutex

	fmt.Println("📡 Downloading mods...")

	for indexSha1, atIndexSha1 := range sha1Hashes {
		// Increment WaitGroup counter
		wg.Add(1)

		go func(index int, sha1 string) {
			// Decrement counter when goroutine completes
			defer wg.Done()

			modName, downloaded, err := downloadViaHash(sha1, version, loader, outputPath)
			if err != nil || !downloaded {
				modName, downloaded, err = downloadViaHash(sha512Hashes[index], version, loader, outputPath)
				if err != nil || !downloaded {
					mu.Lock()
					if modName == "" {
						modName = string(allFiles[index].Name())
					}
					fmt.Printf("❌ Failed: %s for Version: %s / %s\n", modName, version, loader)
					mu.Unlock()
					// Return is used to exit the goroutine
					return
				}
			}

			mu.Lock()
			fmt.Println("✅ Downloaded:", modName)
			mu.Unlock()

		}(indexSha1, atIndexSha1)
	}

	// Wait for all downloads to finish
	wg.Wait()

	fmt.Println("\n\n✅ All downloads completed.")
	scanner.Scan()
}
