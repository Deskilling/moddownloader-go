package modrinth

import (
	"sync"
	"time"

	"moddownloader/extract"
	"moddownloader/filesystem"
	"moddownloader/request"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func Download(id, version, loader, path string) (*extract.Download, error) {
	dl, err := GetDownloads(id, version, loader)
	if err != nil {
		log.Error("failed getting downloads", "id", id, "err", err)
		return nil, err
	}

	log.Debug(dl)

	if dl == nil || dl.Url == "" || dl.Filename == "" {
		log.Error("invalid download object", "id", id)
		return nil, err
	}

	fullPath := path + dl.Filename
	if err = request.DownloadFile(dl.Url, fullPath); err != nil {
		log.Error("failed downloading", "id", id, "url", dl.Url, "err", err)
		return nil, err
	}

	log.Debug("downloaded", "filename", dl.Filename)
	return dl, nil
}

func DownloadAll(id []string, version, loader, output string) (downloadedId []string) {
	filesystem.ClearPath(output)

	if version == "latest" {
		latest, err := request.GetReleaseVersions()
		if err != nil {
			return
		}
		version = latest[0].Version
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, util.GetSettings().General.MaxRoutines)

	rateLimit, err := request.CheckModrinthRate()
	if err != nil {
		log.Error("failed to check rate limit", "err", err)
		return
	}

	remaining := rateLimit.Remaining

	for _, v := range id {
		if remaining <= 0 {
			waitTime := time.Duration(rateLimit.Reset) * time.Second
			log.Warn("rate limit reached, waiting before continuing",
				"wait_seconds", waitTime.Seconds())

			time.Sleep(waitTime)

			rateLimit, err = request.CheckModrinthRate()
			if err != nil {
				log.Error("failed to recheck rate limit", "err", err)
				return
			}

			remaining = rateLimit.Remaining
			log.Info("resumed after cooldown", "remaining", remaining)
		}

		wg.Add(1)
		sem <- struct{}{}
		remaining--

		go func(v string) {
			defer wg.Done()
			defer func() { <-sem }()
			defer func() {
				if r := recover(); r != nil {
					log.Error("panic during download", "id", v, "err", r)
				}
			}()

			_, err := Download(v, version, loader, output)
			if err != nil {
				log.Error("download failed", "id", v, "err", err)
				return
			}

			mu.Lock()
			downloadedId = append(downloadedId, v)
			mu.Unlock()
		}(v)
	}

	wg.Wait()

	return downloadedId
}
