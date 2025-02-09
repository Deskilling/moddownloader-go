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

	fmt.Println("\nPlace all mods into mods_to_update/ and press enter to Continue: ")
	var nichts string
	fmt.Scan(&nichts)

	sha1Hashes, sha512Hashes, err:= calcualteAllHashesFromDirectory("mods_to_update/")
	if err != nil {
		fmt.Println("Error at calcualteAllHashesFromDirectory: ", err)
		return
	}
	
	if len(sha1Hashes) != len(sha512Hashes) {
		fmt.Println("Hash slice have a different size")
		return
	} else {
		fmt.Println("Hash slices are equal")
		lenHashes := len(sha1Hashes)
		fmt.Println(lenHashes)
	}
	// i is the index and v the value at that index
	for i, v := range sha1Hashes {
		fmt.Println(i)
		fmt.Println(v)
	}
}