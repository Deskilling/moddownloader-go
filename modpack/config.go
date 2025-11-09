package modpack

import (
	"fmt"
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

type Modpack struct {
	Name    string   `toml:"Name"`
	Version string   `toml:"Version"`
	Loader  string   `toml:"Loader"`
	Output  string   `toml:"Output"`
	Ids     []string `toml:"Ids"`
}

type Config struct {
	Modpacks map[string]bool `toml:"Modpacks"`
}

func getConfigPath() string {
	return filepath.Join(util.GetSettings().Location.Config, "modpack.toml")
}

func getModpackPath(name string) string {
	return filepath.Join(util.GetSettings().Location.Config, "modpacks", name+".toml")
}

func ReadModpacks() (*Config, error) {
	configPath := getConfigPath()

	if !filesystem.ExistPath(configPath) {
		cfg := &Config{Modpacks: make(map[string]bool)}
		if err := writeConfig(cfg); err != nil {
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
		writeConfig(&cfg)
	}

	if cfg.Modpacks == nil {
		cfg.Modpacks = make(map[string]bool)
	}

	return &cfg, nil
}

func ReadModpack(name string) (*Modpack, error) {
	path := getModpackPath(name)
	if !filesystem.ExistPath(path) {
		return nil, fmt.Errorf("path %s does not exist", path)
	}

	content, err := filesystem.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var mp Modpack
	if err := toml.Unmarshal([]byte(content), &mp); err != nil {
		log.Error("failed unmarshal modpack", "err", err)
		return nil, err
	}

	return &mp, nil
}

func WriteModpackFile(name string, mp Modpack) error {
	data, err := toml.Marshal(&mp)
	if err != nil {
		log.Error("failed marshal modpack", "err", err)
		return err
	}

	path := getModpackPath(name)
	if err := filesystem.CreatePath(filepath.Dir(path)); err != nil {
		log.Error("failed creating modpack folder", "err", err)
		return err
	}

	if err := filesystem.WriteFile(path, data); err != nil {
		log.Error("failed writing modpack file", "err", err)
		return err
	}

	cfg, err := ReadModpacks()
	if err != nil {
		return err
	}

	cfg.Modpacks[name] = true
	return writeConfig(cfg)
}

func RemoveModpack(name string) error {
	path := getModpackPath(name)
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
	return writeConfig(cfg)
}

func writeConfig(cfg *Config) error {
	data, err := toml.Marshal(cfg)
	if err != nil {
		log.Error("failed marshal modpack config", "err", err)
		return err
	}
	return filesystem.WriteFile(getConfigPath(), data)
}

func ConvertMrpack(path string) {
	json := GetMrpackJson(path)
	mp := ParseMrpackJson(json)

	log.Warn(CheckDependencies(*mp))

	tomlmp := Modpack{
		Name:    mp.Name,
		Version: mp.Dependencies["minecraft"],
		Loader:  CheckDependencies(*mp),
		Output:  "./output/" + mp.Name,
		Ids:     GetIdsMrpack(mp),
	}

	WriteModpackFile(mp.Name, tomlmp)
}
