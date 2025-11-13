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
	Config  string
	Modpack string
	Tempdir string
}

type Config struct {
	General   general
	Automatic automatic
	Location  location
}

var dCfg = Config{
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
		Config:  "./moddownloader/",
		Modpack: "./moddownloader/modpack",
		Tempdir: "./moddownloader/temp",
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

	log.Info("Created Default config")
}

// TODO - Add More checks
func checkConfig(cfg *Config) {
	if cfg.General.MaxRoutines <= 0 {
		log.Warn("invalid setting for General.MaxRoutines", "old", cfg.General.MaxRoutines, "new", dCfg.General.MaxRoutines)
		cfg.General.MaxRoutines = dCfg.General.MaxRoutines
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

			checkConfig(&cfg)
			return &cfg
		}
	}

	DefaultConfig()
	cfg = dCfg
	return &cfg
}

func GetSettings() *Config {
	return &cfg
}
