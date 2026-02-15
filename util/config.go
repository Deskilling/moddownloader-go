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


// TODO - Add More checks
func checkConfig(cfg *Config) {
	if cfg.General.MaxRoutines <= 0 {
		log.Warn("invalid setting for General.MaxRoutines", "old", cfg.General.MaxRoutines, "new", dCfg.General.MaxRoutines)
		cfg.General.MaxRoutines = dCfg.General.MaxRoutines
	}
}

func ReadConfig() (*Config, error) {
	if !filesystem.ExistPath("config.toml") {
		w, err := toml.Marshal(&dCfg)
		if err != nil {
			return nil, err
		}

		if err := filesystem.WriteFile("config.toml", w); err != nil {
			log.Error("failed writing default config")
			return nil, err
		}
		log.Info("Created Default config")
		cfg = dCfg
		return &cfg, nil
	}

	c, err := filesystem.ReadFile("config.toml")
	if err != nil {
		log.Error("failed reading file")
		return nil, err
	}

	if len(c) > 0 {
		err = toml.Unmarshal([]byte(c), &cfg)
		if err != nil {
			return nil, err
		}

		checkConfig(&cfg)
		return &cfg, nil
	}

	// if the file is empty, write the default config
	w, err := toml.Marshal(&dCfg)
	if err != nil {
		return nil, err
	}

	if err := filesystem.WriteFile("config.toml", w); err != nil {
		log.Error("failed writing default config")
		return nil, err
	}
	log.Info("Created Default config")
	cfg = dCfg
	return &cfg, nil
}


func GetSettings() *Config {
	return &cfg
}
