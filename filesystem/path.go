package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func ExistPath(path string) bool {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return true
	}
	return false
}

func GetSlug(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

func RemoveMinecraftSuffix(path string) string {
	path = strings.TrimRight(path, `minecraft`)
	return path
}

func ValidPath(path string) string {
	path = strings.TrimRight(path, `/\`)
	path = filepath.Clean(path)
	sep := string(filepath.Separator)
	if !strings.HasSuffix(path, sep) {
		path += sep
	}

	return path
}

func EnsureDirExists(path string) error {
	if ExistPath(path) {
		return nil
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Error("failed creating", "path", path, "err", err)
		return err
	}

	return nil
}

func ClearPath(path, extension string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != extension {
			continue
		}

		if err := os.Remove(fullPath); err != nil {
			return err
		}
	}

	return nil
}

func IsDirEmpty(path string) (bool, error) {
	if !ExistPath(path) {
		return false, os.ErrNotExist
	}

	files, err := os.ReadDir(path)
	if err != nil {
		log.Error("failed reading", "path", path, "err", err)
		return false, err
	}

	return len(files) == 0, nil
}

// verifies if the given path points to a valid .mrpack file
func CheckMrpack(path string) (string, error) {
	modpacksDir := "modpacks"
	pathsToCheck := []string{
		path,
		path + ".mrpack",
		filepath.Join(modpacksDir, path),
		filepath.Join(modpacksDir, path+".mrpack"),
	}

	for _, p := range pathsToCheck {
		stat, err := os.Stat(p)
		if err == nil {
			if stat.IsDir() && filepath.Ext(p) != ".mrpack" {
				continue
			}
			return p, nil
		}
	}

	err := fmt.Errorf("no valid .mrpack file found for %q", path)
	log.Error(err)
	return "", err
}
