package launcher

import (
	"encoding/json"
	"moddownloader/filesystem"
	"moddownloader/request"
	"os"

	"github.com/charmbracelet/log"
)

type cachedRequires struct {
	Suggests string `json:"suggests"`
	Uid      string `json:"uid"`
}

type component struct {
	CachedName     string           `json:"cachedName"`
	CachedRequired []cachedRequires `json:"cachedRequires"`
	CachedVersion  string           `json:"cachedVersion"`
	Important      bool             `json:"important"`
	DependencyOnly bool             `json:"dependencyOnly"`
	Uid            string           `json:"uid"`
	Version        string           `json:"version"`
}

type Mmcpackjson struct {
	Component     []component `json:"components"`
	FormatVersion int         `json:"formatVersion"`
}

// TODO - Check of other launcher also set these env
func IsPrism() bool {
	prismEnv := []string{"INST_NAME", "INST_ID", "INST_DIR", "INST_MC_DIR", "INST_JAVA", "INST_JAVA_ARGS"}

	for _, v := range prismEnv {
		_, b := os.LookupEnv(v)
		if !b {
			continue
		}

		return true
	}

	return false
}

func PrismMmcpack(path string) *Mmcpackjson {
	content, err := filesystem.ReadFile(path)
	if err != nil {
		log.Error("failed reading file", "err", err)
		return nil
	}

	var mmcpack Mmcpackjson
	if err := json.Unmarshal([]byte(content), &mmcpack); err != nil {
		log.Error("failed to unmarshal", "err", err)
		return nil
	}

	return &mmcpack
}

func PrismCurrentVersion() string {
	if !IsPrism() {
		log.Error("prism not detected")
		return ""
	}

	dir := os.Getenv("INST_MC_DIR")
	jsonpath := filesystem.BackOncePath(dir) + "mmc-pack.json"
	mmcpack := PrismMmcpack(jsonpath)

	for _, v := range mmcpack.Component {
		if v.Uid == "net.minecraft" {
			return v.Version
		}
	}

	return ""
}

func PrismCurrentLauncher() string {
	if !IsPrism() {
		log.Error("prism not detected")
		return ""
	}

	dir := os.Getenv("INST_MC_DIR")
	jsonpath := filesystem.BackOncePath(dir) + "mmc-pack.json"
	mmcpack := PrismMmcpack(jsonpath)

	// TODO - Check Names
	for _, v := range mmcpack.Component {
		switch v.Uid {
		case "net.fabricmc.fabric-loader":
			return "fabric"
		case "net.quiltmc":
			return "quilt"
		case "net.forge":
			return "forge"
		case "net.neoforge":
			return "neoforge"
		}
	}

	return ""
}

func PrimsUpdateVersion(version string, mmcpack *Mmcpackjson) error {
	for i, v := range mmcpack.Component {
		switch v.Uid {
		case "net.minecraft", "net.fabricmc.intermediary":
			mmcpack.Component[i].Version = version
		case "net.fabricmc.fabric-loader":
			fabricVer, err := request.GetLatestFabricVersion()
			if err != nil {
				log.Error("failed getting latest fabric version", "err", err)
				return err
			}

			mmcpack.Component[i].Version = fabricVer
		}
		// TODO Add Forge, Neofroge, Quilt and maybe more
	}
	return nil
}

func PrismUpdateJson(version string) {
	if !IsPrism() {
		log.Error("prism not detected")
		return
	}

	if version == "latest" {
		latest, err := request.GetReleaseVersions()
		if err != nil {
			log.Error("failed getting versions", "err", err)
			return
		}
		version = latest[0].Version
	}

	dir := os.Getenv("INST_MC_DIR")
	jsonpath := filesystem.BackOncePath(dir) + "mmc-pack.json"
	mmcpack := PrismMmcpack(jsonpath)

	PrimsUpdateVersion(version, mmcpack)

	mmcjson, err := json.Marshal(mmcpack)
	if err != nil {
		return
	}

	filesystem.WriteFile(jsonpath, mmcjson)
}
