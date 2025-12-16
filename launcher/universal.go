package launcher

import (
	"os"

	"moddownloader/modpack"
	"moddownloader/modrinth"

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

	var id []string
	for _, v := range ids {
		id = append(id, v.Id)
	}

	mp := modpack.Modpack{
		Name:    "default",
		Version: version,
		Loader:  loader,
		Output:  "./mods",
		Ids:     id,
		LastIds: id,
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

	mp.Version = version

	current := modrinth.GetIdsFromPath(path + "/mods")

	currentSet := make(map[string]modrinth.Id, len(current))
	for _, i := range current {
		currentSet[i.Id] = i
	}

	idsSet := make(map[string]struct{}, len(mp.Ids))
	for _, id := range mp.Ids {
		idsSet[id] = struct{}{}
	}

	lastSet := make(map[string]struct{}, len(mp.LastIds))
	for _, id := range mp.LastIds {
		lastSet[id] = struct{}{}
	}

	var add []string
	var remove []string
	var update []string
	var keep []string

	for id := range lastSet {
		if _, ok := currentSet[id]; !ok {
			remove = append(remove, id)
			delete(idsSet, id)
		}
	}

	for id := range currentSet {
		if _, wasBefore := lastSet[id]; !wasBefore {
			add = append(add, id)
			idsSet[id] = struct{}{}
		}
	}

	for id := range idsSet {
		if _, ok := currentSet[id]; ok {
			keep = append(keep, id)
		}
	}

	for id, info := range currentSet {
		if _, ok := idsSet[id]; ok && !info.Latest {
			update = append(update, id)
		}
	}

	mp.Ids = make([]string, 0, len(idsSet))
	for id := range idsSet {
		mp.Ids = append(mp.Ids, id)
	}

	mp.LastIds = currentIdsToList(current)
	modpack.WriteModpackFile("default", *mp)

	for _, id := range keep {
		log.Info("keep", "id", id)
	}
	for _, id := range remove {
		log.Info("remove", "id", id)
	}
	for _, id := range add {
		log.Info("add", "id", id)
	}
	for _, id := range update {
		log.Info("update", "id", id)
	}
}

func currentIdsToList(current []modrinth.Id) []string {
	out := make([]string, len(current))
	for i, j := range current {
		out[i] = j.Id
	}
	return out
}
