package main

import (
	"fmt"
)

func cliMain() {
	err := checkOutputPath("output/")
	if err != nil {
		fmt.Println("Error at checking/creating output/: ", err)
		return
	}

	status, err := doesPathExist("mods_to_update/")
	if err != nil {
		fmt.Println("Error checking/creating mods_to_update/: ", err)
		return
	}

	if status {
		fmt.Println("mods_to_update/ exsits")
	} else {
		fmt.Println("created mods_to_update/")
	}

	// TODO - Change to Buffio scan for defaults
	// Default for version should be the newest
	fmt.Print("Enter Version to use: ")
	var version string
	fmt.Scan(&version)

	fmt.Print("Enter Loader to use: ")
	var loader string
	fmt.Scan(&loader)

	fmt.Println("\nPlace all mods into mods_to_update/ and press enter to Continue: ")
	fmt.Scanln()

	sha1Hashes, sha512Hashes, err:= calcualteAllHashesFromDirectory("mods_to_update/")
	if err != nil {
		fmt.Println("Error at calcualteAllHashesFromDirectory: ", err)
		return
	}

	if len(sha1Hashes) != len(sha512Hashes) {
		fmt.Println("Hash slice have a different size")
		return
	} else {
		lenHashes := len(sha1Hashes)
		fmt.Println(lenHashes)
	}
	// i is the index and v the value at that index
	for indexSha1, atIndexSha1 := range sha1Hashes {
		modName, downloaded, err := downloadViaHash(atIndexSha1, version, loader, "output/")
		if err != nil || !downloaded {
			modName, downloaded, err := downloadViaHash(sha512Hashes[indexSha1], version, loader, "output/")
			if err != nil || !downloaded{
				fmt.Println("Failed to download")
				break 
			} else {
				fmt.Println("Downloaded: ", modName)
			}
		} else {
			fmt.Println("Downloaded: ", modName)
		}
	}
}
