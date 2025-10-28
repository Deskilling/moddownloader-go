package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer modrinthResponse.Body.Close()

	if modrinthResponse.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request failed with status code: %d", modrinthResponse.StatusCode)
	}

	body, err := io.ReadAll(modrinthResponse.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}

func CheckConnection() error {
	_, err := Request(ModrinthEndpoint["default"])
	if err != nil {
		return err
	}
	return nil
}

func DownloadFile(url string, filepath string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", response.StatusCode)
	}

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	return nil
}
