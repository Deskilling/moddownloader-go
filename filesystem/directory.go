package filesystem

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func ReadDirectory(path string, extension string) ([]os.DirEntry, error) {
	if ExistPath(path) {
		if IsDirEmpty(path) {
			log.Error("directory is emptry", "path", path)
			return nil, nil
		}

		allFiles, err := os.ReadDir(path)
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

func CopyDirectory(source string, target string) error {
	err := CreatePath(target)
	if err != nil {
		log.Error("failed to create target directory", "source", source, "target", target, "err", err)
		return err
	}

	files, err := ReadDirectory(source, "")
	if err != nil {
		log.Error("failed to read source directory", "source", source, "target", target, "err", err)
		return err
	}

	for _, f := range files {
		sourcePath := filepath.Join(source, f.Name())
		targetPath := filepath.Join(target, f.Name())

		if f.IsDir() {
			if err := CopyDirectory(sourcePath, targetPath); err != nil {
				return err
			}
		} else {
			content, err := os.ReadFile(sourcePath)
			if err != nil {
				log.Error("failed to read file", "source", source, "err", err)
				return err
			}

			err = os.WriteFile(targetPath, content, os.ModeAppend)
			if err != nil {
				log.Error("failed to write file", "target", target, "err", err)
				return err
			}
		}
	}

	return nil
}
