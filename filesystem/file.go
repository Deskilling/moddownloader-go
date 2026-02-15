package filesystem

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func ReadFile(filepath string) (string, error) {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		log.Error("Error reading file", "filepath", filepath, "err", err)
		return "", err
	}
	return string(fileContent), err
}

func WriteFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		log.Error("failed creating directory", "path", path, "err", err)
		return err
	}

	err := os.WriteFile(path, content, 0644)
	if err != nil {
		log.Error("failed writing", "path", path, "err", err)
		return err
	}
	return nil
}

func DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		// If the file doesn't exist, we don't need to return an error.
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}