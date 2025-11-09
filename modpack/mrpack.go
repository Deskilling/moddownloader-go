package modpack

import (
	"encoding/json"
	"sync"

	"moddownloader/filesystem"
	"moddownloader/modrinth"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

type Mrpack struct {
	Dependencies  map[string]string `json:"dependencies"`
	Files         []file            `json:"files"`
	FormatVersion int               `json:"formatVersion"`
	Game          string            `json:"game"`
	Name          string            `json:"name"`
	VersionId     string            `json:"versionId"`
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
	var idList []string

	sem := make(chan struct{}, util.GetSettings().General.MaxRoutines)
	for _, v := range mrpack.Files {
		wg.Add(1)
		sem <- struct{}{}

		go func(id file) {
			defer wg.Done()
			defer func() { <-sem }()

			hash, err := modrinth.GetIdFromHash(v.Hashes.Sha1)
			if err != nil {
				hash, err = modrinth.GetIdFromHash(v.Hashes.Sha512)
				if err != nil {
					log.Error("failed calulating sha1 and sha512", "err", err)
					return
				}
			}

			idList = append(idList, *hash)
		}(v)
	}

	wg.Wait()

	return idList
}
