package modpack

import (
	"moddownloader/filesystem"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

func writeModpack(mp []Modpack) error {
	w, err := toml.Marshal(mp)
	if err != nil {
		log.Error("failed marshal on mp")
		return err
	}

	filesystem.WriteFile("modpacks.mfg", w)

	return nil
}
