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

func WriteFile(path string, content []byte) {
	err := os.WriteFile(path, content, os.ModeAppend)
	if err != nil {
		return
	}
}
