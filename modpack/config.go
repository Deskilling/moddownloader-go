package modpack

import (

	"path/filepath"

	"moddownloader/filesystem"
	"moddownloader/util"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

type Overwrite struct {
	pathto   string
	pathfrom string
}


type Config struct {
	Modpacks map[string]bool `toml:"Modpacks"`
}

func getConfigPath() string {
	return filepath.Join(util.GetSettings().Location.Config, "modpack.toml")
}

func GetModpackPath(name string) string {
	return filepath.Join(util.GetSettings().Location.Config, "modpacks", name+".toml")
}

func ReadModpacks() (*Config, error) {
	configPath := getConfigPath()
	if err := filesystem.EnsureDirExists(filepath.Dir(configPath)); err != nil {
		log.Error("failed creating modpack config folder", "err", err)
		return nil, err
	}

	if !filesystem.ExistPath(configPath) {
		cfg := &Config{Modpacks: make(map[string]bool)}
		if err := WriteConfig(cfg); err != nil {
			log.Error("failed creating default config", "err", err)
			return nil, err
		}
		return cfg, nil
	}

	content, err := filesystem.ReadFile(configPath)
	if err != nil {
		log.Error("failed reading file", "err", err)
		return nil, err
	}

	var cfg Config
	if err := toml.Unmarshal([]byte(content), &cfg); err != nil {
		log.Warn("invalid modpack config, resetting", "err", err)
		cfg = Config{Modpacks: make(map[string]bool)}
		if err := WriteConfig(&cfg); err != nil {
			log.Error("failed writing new config", "err", err)
			return nil, err
		}
	}

	if cfg.Modpacks == nil {
		cfg.Modpacks = make(map[string]bool)
	}

	return &cfg, nil
}



func RemoveModpack(name string) error {
	path := GetModpackPath(name)
	if filesystem.ExistPath(path) {
		if err := filesystem.DeleteFile(path); err != nil {
			log.Error("failed deleting modpack file", "err", err)
			return err
		}
	}

	cfg, err := ReadModpacks()
	if err != nil {
		return err
	}

	delete(cfg.Modpacks, name)
	return WriteConfig(cfg)
}

func WriteConfig(cfg *Config) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		log.Error("failed marshal modpack config", "err", err)
		return err
	}
	return filesystem.WriteFile(getConfigPath(), data)
}

