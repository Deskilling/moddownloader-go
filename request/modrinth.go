package request

import (
	"encoding/json"
	"fmt"
	"github.com/deskilling/moddownloader-go/extract"
	"io"
	"net/http"
)

func ModrinthWebRequest(endpoint string) (string, error) {
	// Creating a New Request but not sending it
	modrinthRequest, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Setting User Agent for Request
	modrinthRequest.Header.Set("User-Agent", "Deskilling/moddownloader-go")
	// Default http client for the request <- Not needed but why not
	client := &http.Client{}

	// Using the previously created request body to make the request
	modrinthResponse, err := client.Do(modrinthRequest)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	// Closes the Webrequest (important to prevent leaks)
	defer modrinthResponse.Body.Close()

	// Check Response Status
	if modrinthResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", modrinthResponse.StatusCode)
	}

	// Extractin the Information from the Response and checking for errors
	body, err := io.ReadAll(modrinthResponse.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// returning the full requested body if no error happend during the request itself and body
	return string(body), nil
}

func GetReleaseVersions() ([]Version, error) {
	versionsData, err := ModrinthWebRequest(ModrinthEndpoint["availableVersions"])
	if err != nil {
		return nil, fmt.Errorf("error fetching versions: %v", err)
	}

	var versions []Version
	err = json.Unmarshal([]byte(versionsData), &versions)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling versions data: %v", err)
	}

	var releaseVersions []Version
	for _, v := range versions {
		if v.VersionType == "release" {
			releaseVersions = append(releaseVersions, v)
		}
	}

	return releaseVersions, nil
}

func ProjectIdToTitle(projectId string) (string, error) {
	url := fmt.Sprintf(ModrinthEndpoint["modInformation"], projectId)
	response, err := ModrinthWebRequest(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch project information: %w", err)
	}

	extractedInformation, err := extract.Mod(response)
	if err != nil {
		return "", fmt.Errorf("failed to parse project information: %w", err)
	}

	return extractedInformation.ProjectTitle, nil
}
