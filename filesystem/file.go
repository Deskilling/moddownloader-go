package filesystem

import (
	"os"

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
	if !ExistPath(path) {
		CreatePath(path)
	}

	err := os.WriteFile(path, content, 0777)
	if err != nil {
		log.Error("failed writing", "path", path, "err", err)
		return err
	}
	return nil
}
