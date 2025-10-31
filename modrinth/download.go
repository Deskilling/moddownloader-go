package modrinth

import (
	"sync"

	"moddownloader/extract"
	"moddownloader/request"
	"moddownloader/util"

	"github.com/charmbracelet/log"
)

func Download(id, version, loader, path string) (*extract.Download, error) {
	dl, err := GetDownloads(id, version, loader)
	if err != nil {
		log.Error("failed getting downloads", "err", err)
		return nil, err
	}

	if err = request.DownloadFile(dl.Url, path+dl.Filename); err != nil {
		log.Error("failed downloading", "url", dl.Url, "err", err)
	}

	return dl, nil
}

func DownloadAll(id []string, version, loader, output string) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, util.GetSettings().General.MaxRoutines)
	for _, v := range id {
		wg.Add(1)
		sem <- struct{}{}

		go func(id string) {
			defer wg.Done()
			defer func() { <-sem }()

			Download(id, version, loader, output)
		}(v)
	}

	wg.Wait()
}
