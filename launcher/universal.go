package launcher

import (
	"moddownloader/filesystem"
	"moddownloader/modpack"
	"moddownloader/modrinth"
	"os"

	"github.com/charmbracelet/log"
)

// Execution dir is expected at minecraft instance
func CreateModpack(version, loader string) {
	path, err := os.Getwd()
	if err != nil {
		log.Error("failed getting execution path", "err", err)
		return
	}

	log.Debug(path)
	ids := modrinth.GetIdsFromPath(path + "/mods")

	mp := modpack.Modpack{
		Name:    "default",
		Version: version,
		Loader:  loader,
		Output:  "./mods",
		Ids:     ids,
		LastIds: ids,
	}

	modpack.WriteModpackFile("default", mp)
}

func ReadModpack() (*modpack.Modpack, error) {
	mp, err := modpack.ReadModpack("default")
	if err != nil {
		log.Error("failed reading modpack file", "err", err)
		return nil, err
	}

	return mp, nil
}

func UpdateModpack(version string, mp *modpack.Modpack) {
	path, err := os.Getwd()
	if err != nil {
		log.Error("failed getting execution path", "err", err)
		return
	}

	currentIds := modrinth.GetIdsFromPath(path + "/mods")

	if len(currentIds) == 0 {
		log.Info("no existing mod IDs found, downloading all mods")
		downloaded := modrinth.DownloadAll(mp.Ids, version, mp.Loader, filesystem.ValidPath(mp.Output))
		mp.LastIds = downloaded
		return
	}

	mapOverview := make(map[string]bool)
	mapOld := make(map[string]bool)
	mapNew := make(map[string]bool)

	for _, id := range mp.Ids {
		mapOverview[id] = true
	}
	for _, id := range mp.LastIds {
		mapOld[id] = true
	}
	for _, id := range currentIds {
		mapNew[id] = true
	}

	for id := range mapOld {
		if !mapNew[id] {
			delete(mapOverview, id)
		}
	}

	for id := range mapNew {
		if !mapOld[id] {
			mapOverview[id] = true
		}
	}

	result := make([]string, 0, len(mapOverview))
	for id := range mapOverview {
		result = append(result, id)
	}

	mp.Ids = result

	log.Debug("loader is", "loader", mp.Loader)
	log.Debug("version is", "version", version, "prismversion", PrismCurrentVersion())

	PrismUpdateJson(version)

	downloaded := modrinth.DownloadAll(mp.Ids, version, mp.Loader, filesystem.ValidPath(mp.Output))
	mp.LastIds = downloaded
}
