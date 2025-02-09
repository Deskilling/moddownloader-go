package main

// ---- This Includes all the Code used for GET request in the Downloader ---- //

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

	// "search": "https://api.modrinth.com/v2/search",
}

// TODO -- Implement a header for Request for
func modrinthWebRequest(endpoint string) (string, error) {
	// Creating a New Request but not sending it
	modrinthRequest, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error Creating Request:", err)
		return "", err
	}

	// Setting User Agent for Request
	modrinthRequest.Header.Set("User-Agent", "Deskilling/moddownloader-go")

	// Default http client for the request <- Not needed but why not
	client := &http.Client{}
	// Using the previously created request body to make the request
	modrinthResponse, err := client.Do(modrinthRequest)
	if err != nil {
		fmt.Println("Error at Response: ", err)
		return "", err
	}

	// Closes the Webrequest (important to prevent leaks)
	defer modrinthResponse.Body.Close()

	// Extractin the Information from the Response and checking for errors
	body, err := io.ReadAll(modrinthResponse.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return "", err
	}

	// returning the full requested body if no error happend during the request itself and body
	// TODO - should check what kind of erors can hapen (the most common ones. No Internet or Invalid request or missing mod and stuff)
	return string(body), nil
}

func downloadFromUrl(url string, filepath string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
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

	// Copies the Chunks from the Request and Writes it (32kb at a time)
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func projectIdToTitle(projectId string) (string, error) {
	url := fmt.Sprintf(modrinthEndpoint["modInformation"], projectId)
	response, err := modrinthWebRequest(url)
	if err != nil {
		return "", err
	}

	extractedInformation, err := extractModInformation(response)
	if err != nil {
		return "", err
	}

	return extractedInformation.ProjectTitle, nil
}
