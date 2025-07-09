package request

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filepath string) (err error) {
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
