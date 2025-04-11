package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// The map for the different endpoints <- probabbly useless but it looks cool
type EndpointMap map[string]string

// All for this project releveant Endpoints (or that might be added in the future)
var modrinthEndpoint = EndpointMap{
	"default":               "https://api.modrinth.com",
	"modInformation":        "https://api.modrinth.com/v2/project/%s",
	"modVersionInformation": "https://api.modrinth.com/v2/project/%s/version",
	"versionFileHash":       "https://api.modrinth.com/v2/version_file/%s",
	"versionUpdate":         "https://api.modrinth.com/v2/version_file/{hash}/update",
	"availableVersions":     "https://api.modrinth.com/v2/tag/game_version",
	"availableLoaders":      "https://api.modrinth.com/v2/tag/loader",
	// GRRR
	"fabricVersions": "https://meta.fabricmc.net/v2/versions/loader",

	// "search": "https://api.modrinth.com/v2/search",
}

func modrinthWebRequest(endpoint string) (string, error) {
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

func downloadFile(url string, filepath string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	// Closes File
	defer out.Close()

	// Get the data
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	// Closes The Request
	defer response.Body.Close()

	// check status code
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", response.StatusCode)
	}

	// Copies the Chunks from the Request and Writes it (32kb at a time)
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}
