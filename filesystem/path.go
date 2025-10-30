package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func ExistPath(path string) bool {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true
	}
	return false
}

func CreatePath(path string) error {
	if ExistPath(path) {
		return nil
	}

	err := os.MkdirAll(path, os.ModeDir)
	if err != nil {
		log.Error("failed creating", "path", path, "err", err)
		return err
	}

	return nil
}

func IsDirEmpty(path string) bool {
	if !ExistPath(path) {
		return true
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Error("failed reading", "path", path, "err", err)
		return true
	}

	if len(files) != 0 {
		return false
	}

	return true
}

// verifies if the given path points to a valid .mrpack file
func CheckMrpack(path string) (string, error) {
	modpacksDir := "modpacks"
	pathsToCheck := []string{
		filepath.Join(path),                        // Default
		filepath.Join(path + ".mrpack"),            // path.mrpack
		filepath.Join(modpacksDir, path),           // modpacks/path
		filepath.Join(modpacksDir, path+".mrpack"), // modpacks/path.mrpack
	}

	var validPath string
	for _, p := range pathsToCheck {
		if stat, err := os.Stat(p); err == nil {
			if stat.IsDir() {
				if filepath.Ext(p) == ".mrpack" {
					validPath = p
					break
				}
			} else {
				validPath = p
				break
			}
		}
	}

	if validPath == "" {
		err := fmt.Errorf("invalid path")
		log.Error("no valid .mrpack file found", "path", path, "err", err)
		return "", err
	}

	if filepath.Ext(validPath) != ".mrpack" {
		err := fmt.Errorf("invalid extension")
		log.Error("returned file does not contain .mrpack extension", "path", path)
		return "", err
	}

	return validPath, nil
}
