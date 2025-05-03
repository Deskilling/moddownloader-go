package main

import "fmt"

func main() {
	downloadALlModpack("modpacks/EumelcraftPack.mrpack", "1.21.5", "fabric")
	fmt.Println(getLatestQuiltVersion())
	file, _ := getLauncherProfiles()
	ewaoiurhgn := parseLauncherProfiles(file)

	profileAdd(ewaoiurhgn, "fabric", "1.21.5", "schwanzred")

	checkOutputPath("test/temp")
	exportModpack("modpacks/EumelcraftPack.mrpack", "test/temp")

	/*
			err := checkConnection()
			if err != nil {
				return
			}

			if len(os.Args) < 2 {
				cliMain()
			} else {
				runArgs()
			}
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("\n[Enter to exit]")
		scanner.Scan()

	*/
}
