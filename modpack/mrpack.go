package modpack

import (
	"encoding/json"
	"sync"

	"moddownloader/filesystem"
	"moddownloader/modrinth"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func ExtractMrPack(path string) {
	filesystem.ExtractZip(path, util.GetSettings().Location.Tempdir+"modpack"+path)

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

			log.Debug("calculated", "idlist", idList)
			idList = append(idList, *hash)
		}(v)
	}

	wg.Wait()

	return idList
}
