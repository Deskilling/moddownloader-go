package util

import (
	"encoding/json"

	"github.com/deskilling/moddownloader-go/filesystem"
	"github.com/deskilling/moddownloader-go/request"
)

type Config struct {
	Mode    string
	Version string
	Loader  string
	Input   string
	Output  string
}

func CreateConfig() {
	var config Config
	latest, err := request.GetReleaseVersions()
	if err != nil {
		return
	}

	config.Mode = "mods"
	config.Version = latest[0].Version
	config.Loader = "fabric"
	config.Input = "./input"
	config.Output = "./output"

	cfgJson, err := json.Marshal(config)
	if err != nil {
		return
	}

	// TODO Pretty json
	filesystem.WriteFile("mod.cfg", []byte(string(cfgJson)))
}
func LoadConfig() (Config, error) {
	configJson, err := filesystem.ReadFile("mod.cfg")

	if err != nil || configJson == "" {
		return Config{}, err
	}

	var cfg Config
	json.Unmarshal([]byte(configJson), &cfg)

	return cfg, nil
}

func GetEmptyConfig() Config {
	var config Config
	return config
}

func SaveConfig() {

}
