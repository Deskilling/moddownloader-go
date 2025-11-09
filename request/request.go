package request

import (
	"fmt"
	"io"
	"moddownloader/filesystem"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
)

type EndpointMap map[string]string

var ModrinthEndpoint = EndpointMap{
	"default":               "https://api.modrinth.com",
	"modInformation":        "https://api.modrinth.com/v2/project/%s",
	"modVersionInformation": "https://api.modrinth.com/v2/project/%s/version",
	"versionFileHash":       "https://api.modrinth.com/v2/version_file/%s",
	"versionUpdate":         "https://api.modrinth.com/v2/version_file/{hash}/update",
	"availableVersions":     "https://api.modrinth.com/v2/tag/game_version",
	"availableLoaders":      "https://api.modrinth.com/v2/tag/loader",

	// "search": "https://api.modrinth.com/v2/search",
}

func Request(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Deskilling/moddownloader-go")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("request failed", "err", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		log.Error("request failed", "status", resp.StatusCode)
		return nil, fmt.Errorf("status not ok")
	}

	return resp, nil
}

func GetBody(endpoint string) (string, error) {
	resp, err := Request(endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read response body", "err", err)
		return "", err
	}

	return string(body), nil
}

type ModrinthRateLimit struct {
	Limit     int
	Remaining int
	Reset     int
}

func CheckModrinthRate() (rateLimit ModrinthRateLimit, err error) {
	response, err := Request(ModrinthEndpoint["default"])
	if err != nil {
		return rateLimit, err
	}
	defer response.Body.Close()

	rateLimit.Limit, err = parse(response, "X-RateLimit-Limit")
	rateLimit.Remaining, err = parse(response, "X-RateLimit-Remaining")
	rateLimit.Reset, err = parse(response, "X-RateLimit-Reset")

	return rateLimit, err
}

func parse(r *http.Response, name string) (int, error) {
	s := r.Header.Get(name)
	if s == "" {
		return 0, fmt.Errorf("missing header %s", name)
	}
	v, e := strconv.ParseInt(s, 10, 0)
	if e != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, e)
	}
	return int(v), nil
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
