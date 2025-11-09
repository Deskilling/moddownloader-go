package modrinth

import (
	"fmt"

	"moddownloader/extract"
	"moddownloader/request"

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

func GetIdFromHash(hash string) (*string, error) {
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

	return &emvi.ProjectId, nil
}
