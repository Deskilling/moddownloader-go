package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/user"
	"runtime"
)

type Settings struct {
	EnableAnalytics  bool   `json:"enableAnalytics"`
	EnableAdvanced   bool   `json:"enableAdvanced"`
	KeepLauncherOpen bool   `json:"keepLauncherOpen"`
	SoundOn          bool   `json:"soundOn"`
	ShowMenu         bool   `json:"showMenu"`
	EnableSnapshots  bool   `json:"enableSnapshots"`
	EnableHistorical bool   `json:"enableHistorical"`
	EnableReleases   bool   `json:"enableReleases"`
	ProfileSorting   string `json:"profileSorting"`
	ShowGameLog      bool   `json:"showGameLog"`
	CrashAssistance  bool   `json:"crashAssistance"`
}

type Profile struct {
	ID            string `json:"id"`
	LastUsed      string `json:"lastUsed"`
	LastVersionId string `json:"lastVersionId"`
	Created       string `json:"created"`
	Icon          string `json:"icon"`
	Name          string `json:"name"`
	Type          string `json:"type"`
}

// Config represents the entire JSON structure
type Config struct {
	Settings Settings           `json:"settings"`
	Profiles map[string]Profile `json:"profiles"`
	Version  int                `json:"version"`
}

func getLauncherProfiles() string {
	var path string

	username, _ := user.Current()

	operatingSystem := runtime.GOOS
	if operatingSystem == "windows" {
		path = fmt.Sprintf("%s\\AppData\\Roaming\\.minecraft\\launcher_profiles.json", username.HomeDir)

	} else if operatingSystem == "linux" {
		path = fmt.Sprintf("%s/.minecraft/launcher_profiles.json", username.HomeDir)

	} else if operatingSystem == "darwin" {
		// not working bruh <- i think
		path = fmt.Sprintf("%s/Library/Application Support/minecraft/launcher_profiles.json", username.HomeDir)
	}

	fmt.Println(path)

	return readFile(path)
}

func parseLauncherProfiles(jsonData string) {
	var config Config
	err := json.Unmarshal([]byte(jsonData), &config)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	for _, profile := range config.Profiles {
		fmt.Println(profile.Type)
	}

}
