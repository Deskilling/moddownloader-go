package modrinth

import (
	"fmt"

	"moddownloader/extract"
	"moddownloader/request"

	"github.com/charmbracelet/log"
)

func ProjectIdToTitle(projectId string) (string, error) {
	url := fmt.Sprintf(request.ModrinthEndpoint["modInformation"], projectId)
	response, err := request.Request(url)
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
