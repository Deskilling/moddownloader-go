package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(filepath string) string {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}
	return string(fileContent)
}

func WriteFile(path string, content []byte) {
	// 0064 is the permission level
	err := os.WriteFile(path, content, 0064)
	if err != nil {
		return
	}
}

func GetAllFilesFromDirectory(directory string, extension string) ([]os.DirEntry, error) {
	doesExist, err := DoesPathExist(directory)
	if err != nil {
		return nil, err
	}
	if doesExist {
		allFiles, err := os.ReadDir(directory)
		if err != nil {
			return nil, err
		}

		var filteredFiles []os.DirEntry

		for _, file := range allFiles {
			if filepath.Ext(file.Name()) == extension {
				filteredFiles = append(filteredFiles, file)
			}
		}

		return filteredFiles, nil
	}
	return nil, nil
}
