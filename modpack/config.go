package modpack

import (
	"moddownloader/filesystem"
	"moddownloader/util"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

func ReadModpacks() (*Config, error) {
	if filesystem.ExistPath(util.GetSettings().Location.Modpacks) {
		c, err := filesystem.ReadFile(util.GetSettings().Location.Modpacks)
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
	}

	return nil, nil
}

func AddModpack(name, version, loader, input, output string) error {
	mp, err := ReadModpacks()
	if err != nil {
		log.Error("failed reading modpacks", "err", err)
		return err
	}

	nmp := Modpack{
		Version: version,
		Loader:  loader,
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

	if err := filesystem.WriteFile(util.GetSettings().Location.Modpacks, w); err != nil {
		log.Error("failed writing file", "err", err)
		return err
	}

	return nil
}
