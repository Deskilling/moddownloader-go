package profiles

import (
	"encoding/json"
	"fmt"
	"os/user"
	"runtime"

	"github.com/deskilling/moddownloader-go/filesystem"
	"github.com/deskilling/moddownloader-go/util"
)

func getLauncherProfiles() (string, string) {
	var path string

	username, _ := user.Current()

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		path = fmt.Sprintf("%s\\AppData\\Roaming\\.minecraft\\launcher_profiles.json", username.HomeDir)

	case "linux":
		path = fmt.Sprintf("%s/.minecraft/launcher_profiles.json", username.HomeDir)

	case "darwin":
		// TODO - Test
		path = fmt.Sprintf("%s/Library/Application Support/minecraft/launcher_profiles.json", username.HomeDir)
	}

	return filesystem.ReadFile(path), path
}

func parseLauncherProfiles(jsonData string) Config {
	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		panic(err)
	}
	return config
}

func profileAdd(parsedJson Config, loader string, version string, name string) {
	profileId := latestProfile(loader, version)

	profile, exists := parsedJson.Profiles[profileId]
	if exists {
		profileId = fmt.Sprintf("%s-%s-%s", profileId, name, version)
	}

	profile.Name = fmt.Sprintf("%s %s", name, version)
	profile.LastVersionId = profileLastVersionId(loader, version)

	currentTime := util.GetTime()
	profile.LastUsed = currentTime
	profile.Created = currentTime

	profile.Icon = profileIcon(loader)

	parsedJson.Profiles[profileId] = profile

	_, path := getLauncherProfiles()

	jsonOutput, err := json.MarshalIndent(parsedJson, "", "  ")
	if err != nil {
		fmt.Println("❌ Error marshaling JSON:", err)
		return
	}

	filesystem.WriteFile(path, jsonOutput)
}
