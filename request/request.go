package request

import (
	"fmt"
	"io"
	"moddownloader/filesystem"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
)

func Request(endpoint string) (string, error) {
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	request.Header.Set("User-Agent", "Deskilling/moddownloader-go")
	client := &http.Client{}

	modrinthResponse, err := client.Do(request)
	if err != nil {
		log.Error("request failed", "err", err)
		return "", err
	}
	defer modrinthResponse.Body.Close()

	if modrinthResponse.StatusCode != http.StatusOK {
		log.Error("request failed", "status", modrinthResponse.StatusCode)
		return "", fmt.Errorf("status not ok")
	}

	body, err := io.ReadAll(modrinthResponse.Body)
	if err != nil {
		log.Error("failed to read response body", "err", err)
		return "", err
	}
	return string(body), nil
}

func CheckConnection() bool {
	_, err := Request(ModrinthEndpoint["default"])
	if err != nil {
		return false
	}
	return true
}

func DownloadFile(url string, filepath string) (err error) {
	if err := filesystem.CreatePath(filepath); err != nil {
		log.Error("failed creating path", "err", err)
		return err
	}

	f, err := os.Create(filepath)
	if err != nil {
		log.Error("failed to create file", "err", err)
		return err
	}
	defer f.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Error("failed to download file", "status", response.StatusCode)
		return fmt.Errorf("status is not ok")
	}

	_, err = io.Copy(f, response.Body)
	if err != nil {
		log.Error("failed to download file", "err", err)
		return err
	}

	return nil
}
