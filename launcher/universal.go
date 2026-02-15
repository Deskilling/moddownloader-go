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
	currentSet := createIdSet(current)
	idsSet := createStringSet(mp.Ids)
	lastSet := createStringSet(mp.LastIds)

	remove := getRemovedMods(lastSet, currentSet)
	for _, id := range remove {
		delete(idsSet, id)
	}

	add := getAddedMods(lastSet, currentSet)
	for _, id := range add {
		idsSet[id] = struct{}{}
	}

	update := getUpdateMods(idsSet, currentSet)

	mp.Ids = stringSetToSlice(idsSet)
	mp.LastIds = currentIdsToList(current)
	modpack.WriteModpackFile("default", *mp)

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

func getRemovedMods(lastSet map[string]struct{}, currentSet map[string]modrinth.Id) []string {
	var remove []string
	for id := range lastSet {
		if _, ok := currentSet[id]; !ok {
			remove = append(remove, id)
		}
	}
	return remove
}

func getAddedMods(lastSet map[string]struct{}, currentSet map[string]modrinth.Id) []string {
	var add []string
	for id := range currentSet {
		if _, wasBefore := lastSet[id]; !wasBefore {
			add = append(add, id)
		}
	}
	return add
}

func getUpdateMods(idsSet map[string]struct{}, currentSet map[string]modrinth.Id) []string {
	var update []string
	for id, info := range currentSet {
		if _, ok := idsSet[id]; ok && !info.Latest {
			update = append(update, id)
		}
	}
	return update
}

func createIdSet(ids []modrinth.Id) map[string]modrinth.Id {
	set := make(map[string]modrinth.Id, len(ids))
	for _, i := range ids {
		set[i.Id] = i
	}
	return set
}

func createStringSet(ids []string) map[string]struct{} {
	set := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}

func stringSetToSlice(set map[string]struct{}) []string {
	slice := make([]string, 0, len(set))
	for id := range set {
		slice = append(slice, id)
	}
	return slice
}

func currentIdsToList(current []modrinth.Id) []string {
	out := make([]string, len(current))
	for i, j := range current {
		out[i] = j.Id
	}
	return out
}
