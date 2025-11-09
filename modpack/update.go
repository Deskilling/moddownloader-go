package modpack

import (
	"moddownloader/modrinth"

	"github.com/charmbracelet/log"
)

func UpdateToml(name, version string) error {
	mp, err := ReadModpack(name)
	if err != nil {
		log.Error("toml not found", "err", err)
		return err
	}

	log.Debug("Modpack:\n", "version", version, "loader", mp.Loader, "output", mp.Output, "ids", mp.Ids)

	modrinth.DownloadAll(mp.Ids, version, mp.Loader, mp.Output)

	return nil
}
