package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func CopyDirectory(source string, target string) error {
	err := os.MkdirAll(target, 0755) // for dir
	if err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	files, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
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
				return fmt.Errorf("failed to read file %s: %w", sourcePath, err)
			}

			err = os.WriteFile(targetPath, content, 0644) // for files
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
		}
	}

	return nil
}
