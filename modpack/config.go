package modpack

import (
	"moddownloader/filesystem"
	"moddownloader/util"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

type Modpack struct {
	Version string
	Loader  string
	Mode    string // Either mrpack or toml
	Input   string // Either mrpack or toml
	Output  string
}

type Config struct {
	Modpacks map[string]Modpack
}

func ReadModpacks() (*Config, error) {
	if filesystem.ExistPath(util.GetSettings().Location.Modpack) {
		c, err := filesystem.ReadFile(util.GetSettings().Location.Config + "modpack.toml")
		if err != nil {
			log.Error("failed reading file")
			return nil, err
		}

		if len(c) > 0 {
			var mp Config
			err = toml.Unmarshal([]byte(c), &mp)
			if err != nil {
				log.Error("failed unmarshal", "err", err)
				return nil, err
			}
			return &mp, nil
		}
	} else {
		if err := filesystem.CreatePath(util.GetSettings().Location.Config + "modpack.toml"); err != nil {
			log.Error("failed creating modpack config", "err", err)
			return nil, err
		}

		// This may be used in the future to create default settings

		var config Config
		w, err := toml.Marshal(config)
		if err != nil {
			return nil, err
		}

		if filesystem.WriteFile(util.GetSettings().Location.Config+"modpack.toml", w) != nil {
			log.Error("failed writing default config")
			return nil, err
		}
	}

	return nil, nil
}

func AddModpack(name, version, loader, mode, input, output string) error {
	mp, err := ReadModpacks()
	if err != nil {
		log.Error("failed reading modpacks", "err", err)
		return err
	}

	nmp := Modpack{
		Version: version,
		Loader:  loader,
		Mode:    mode,
		Input:   input,
		Output:  output,
	}

	mp.Modpacks[name] = nmp

	WriteModpack(*mp)

	return nil
}

func RemoveModpack(name string) error {
	mp, err := ReadModpacks()
	if err != nil {
		log.Error("failed reading modpacks", "err", err)
		return err
	}

	delete(mp.Modpacks, name)

	if err := WriteModpack(*mp); err != nil {
		log.Error("failed updating modpack", "err", err)
		return err
	}

	return nil
}

func WriteModpack(mp Config) error {
	w, err := toml.Marshal(&mp)
	if err != nil {
		log.Error("failed marshal on mp")
		return err
	}

	if err := filesystem.WriteFile(util.GetSettings().Location.Config+"modpack.toml", w); err != nil {
		log.Error("failed writing file", "err", err)
		return err
	}

	return nil
}

func ConvertMrpack(path string) {
}
