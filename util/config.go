package util

import (
	"moddownloader/filesystem"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

var cfg Config

type general struct {
	MaxRoutines int `comment:"maximum of gorutines at once"`
}

type automatic struct {
	Toggle   bool
	Modpacks map[string]string `comment:"modpacks to update (name = version)"`
}

type location struct {
	Modpacks string `comment:"path to the modpack.toml file"`
}

type Config struct {
	General   general
	Automatic automatic
	Location  location
}

var dCfg = &Config{
	General: general{
		MaxRoutines: 64,
	},
	Automatic: automatic{
		Toggle: false,
		Modpacks: map[string]string{
			"eumelpack":    "latest",
			"speedrunning": "1.16.1",
		},
	},
	Location: location{
		Modpacks: "./config/modpack.toml",
	},
}

func DefaultConfig() {
	w, err := toml.Marshal(&dCfg)
	if err != nil {
		return
	}

	log.Debug("", "w", string(w))

	if filesystem.WriteFile("config.toml", w) != nil {
		log.Error("failed writing default config")
		return
	}
}

func ReadConfig() *Config {
	if filesystem.ExistPath("config.toml") {
		c, err := filesystem.ReadFile("config.toml")
		if err != nil {
			log.Error("failed reading file")
			return nil
		}

		if len(c) > 0 {
			err = toml.Unmarshal([]byte(c), &cfg)
			if err != nil {
				return nil
			}
			return &cfg
		}
	}

	DefaultConfig()
	return dCfg
}

func GetSettings() *Config {
	return &cfg
}
