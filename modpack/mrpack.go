package modpack

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"moddownloader/filesystem"
	"moddownloader/modrinth"
	"moddownloader/request"
	"moddownloader/util"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

type Mrpack struct {
	Dependencies  map[string]string `json:"dependencies"`
	Files         []file            `json:"files"`
	FormatVersion int               `json:"formatVersion"`
	Game          string            `json:"game"`
	Name          string            `json:"name"`
	VersionId     string            `json:"versionId"`
}

type Modpack struct {
	Name    string   `toml:"Name"`
	Version string   `toml:"Version"`
	Loader  string   `toml:"Loader"`
	Output  string   `toml:"Output"`
	Ids     []string `toml:"Ids"`
	LastIds []string `tomem:"Lastds"`
}

type file struct {
	DownloadUrl []string          `json:"downloads"`
	Env         map[string]string `json:"env"`
	FileSize    int               `json:"fileSize"`
	Hashes      hashes            `json:"hashes"`
	Path        string            `json:"path"`
}

type hashes struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
}

func GetMrpackJson(path string) (content string) {
	modpackPath := util.GetSettings().Location.Tempdir + filesystem.GetSlug(path) + "/"
	log.Debug("extracting to", "path", modpackPath)
	err := filesystem.ExtractZip(path, modpackPath)
	if err != nil {
		log.Error("extracting mrpack", "err", err)
		return ""
	}
	content, err = filesystem.ReadFile(modpackPath + "modrinth.index.json")
	if err != nil {
		log.Error("reading json", "err", err)
		return ""
	}

	return content
}

func ParseMrpackJson(jsonData string) *Mrpack {
	var mrpack Mrpack
	if err := json.Unmarshal([]byte(jsonData), &mrpack); err != nil {
		log.Error("failed umarshal", "err", err)
	}

	return &mrpack
}

func GetIdsMrpack(mrpack *Mrpack) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var idList []string

	sem := make(chan struct{}, util.GetSettings().General.MaxRoutines)

	rateLimit, err := request.CheckModrinthRate()
	if err != nil {
		log.Error("failed to check rate limit", "err", err)
		return nil
	}
	remaining := rateLimit.Remaining

	for _, v := range mrpack.Files {
		if remaining <= 0 {
			waitTime := time.Duration(rateLimit.Reset) * time.Second
			log.Warn("rate limit reached, waiting before continuing",
				"wait_seconds", waitTime.Seconds())

			time.Sleep(waitTime)

			rateLimit, err = request.CheckModrinthRate()
			if err != nil {
				log.Error("failed to recheck rate limit", "err", err)
				return nil
			}
			remaining = rateLimit.Remaining
			log.Info("resumed after cooldown", "remaining", remaining)
		}

		wg.Add(1)
		sem <- struct{}{}
		remaining--

		go func(file file) {
			defer wg.Done()
			defer func() { <-sem }()

			defer func() {
				if r := recover(); r != nil {
					log.Error("panic during GetIdFromHash", "file", file, "err", r)
				}
			}()

			hash, err := modrinth.GetIdFromHash(file.Hashes.Sha1)
			if err != nil {
				hash, err = modrinth.GetIdFromHash(file.Hashes.Sha512)
				if err != nil {
					log.Error("failed calculating sha1 and sha512", "err", err)
					return
				}
			}

			mu.Lock()
			idList = append(idList, hash.ProjectId)
			mu.Unlock()
		}(v)
	}

	wg.Wait()
	return idList
}

func CheckDependencies(mp Mrpack) string {
	loaders, err := request.GetAllLoaders()
	if err != nil {
		return ""
	}

	// only 4 types of modpacks are currently on modrinth
	// fabric and quilt both have -loder added
	// forge and neoforge just use their name
	for _, loader := range loaders {
		if loader.Name == "minecraft" {
			continue
		}
		_, ok := mp.Dependencies[loader.Name]
		if ok {
			return loader.Name
		}

		_, ok = mp.Dependencies[loader.Name+"-loader"]
		if ok {
			return loader.Name
		}
	}

	return ""
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
		LastIds: GetIdsMrpack(mp),
	}

	WriteModpackFile(mp.Name, tomlmp)
}
func WriteModpackFile(name string, mp Modpack) error {
	data, err := toml.Marshal(&mp)
	if err != nil {
		log.Error("failed marshal modpack", "err", err)
		return err
	}

	path := GetModpackPath(name)
	if err := filesystem.EnsureDirExists(filepath.Dir(path)); err != nil {
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
	return WriteConfig(cfg)
}

func ReadModpack(name string) (*Modpack, error) {
	path := GetModpackPath(name)
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