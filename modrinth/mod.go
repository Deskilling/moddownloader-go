package modrinth

import (
	"fmt"
	"sync"
	"time"

	"moddownloader/extract"
	"moddownloader/filesystem"
	"moddownloader/request"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func GetDownloads(id, version, loader string) (*extract.Download, error) {
	cInfo, err := request.GetBody(fmt.Sprintf(request.ModrinthEndpoint["modVersionInformation"], id))
	if err != nil {
		log.Error("request failed", "err", err)
		return nil, err
	}

	mInfo, err := extract.AllVersionHash(cInfo)
	if err != nil {
		log.Error("failed extracting all versions", "err", err)
		return nil, err
	}

	dl, err := extract.GetDownload(*mInfo, version, loader)
	if err != nil {
		log.Error("failed extracting downloads", "err", err)
		return nil, err
	}

	return dl, nil
}

func ProjectIdToTitle(projectId string) (string, error) {
	url := fmt.Sprintf(request.ModrinthEndpoint["modInformation"], projectId)
	response, err := request.GetBody(url)
	if err != nil {
		log.Error("failed to fetch project information", "err", err)
		return "", err
	}

	extractedInformation, err := extract.Mod(response)
	if err != nil {
		log.Error("failed to parse project information", "err", err)
		return "", err
	}

	return extractedInformation.ProjectTitle, nil
}

func GetIdFromHash(hash string) (*extract.ModVersionInformation, error) {
	url := fmt.Sprintf(request.ModrinthEndpoint["versionFileHash"], hash)
	response, err := request.GetBody(url)
	if err != nil {
		log.Error("failed to fetch project information", "url", url, "err", err)
		return nil, err
	}

	emvi, err := extract.Version(response)
	if err != nil {
		log.Error("failed extracting allversions", "hash", hash, "err", err)
	}

	return emvi, nil
}

type Id struct {
	Latest bool
	Id     string
}

func GetIdsFromPath(path string) (ids []Id) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	sem := make(chan struct{}, util.GetSettings().General.MaxRoutines)

	rateLimit, err := request.CheckModrinthRate()
	if err != nil {
		log.Error("failed to check rate limit", "err", err)
		return nil
	}
	remaining := rateLimit.Remaining

	sha1, sha512, files, err := filesystem.CalculateDirectoryHashes(path, ".jar")
	if err != nil {
		log.Error("failed caluclating hashes", "err", err)
		return
	}

	for i := range files {
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

		go func(sha1, sha512 string) {
			defer wg.Done()
			defer func() { <-sem }()

			hashId, err := GetIdFromHash(sha1)
			if err != nil {
				hashId, err = GetIdFromHash(sha512)
				if err != nil {
					log.Warn("Not found on Modrinth (skipping)", "filename", files[i])
					return
				}
			}

			var new bool = false
			if sha1 == hashId.Files[0].Hashes.Sha1 && sha512 == hashId.Files[0].Hashes.Sha512 {
				new = true
			}

			mu.Lock()
			var rId Id = Id{
				Latest: new,
				Id:     hashId.ProjectId,
			}
			ids = append(ids, rId)
			mu.Unlock()

		}(sha1[i], sha512[i])
	}

	wg.Wait()
	return ids
}
