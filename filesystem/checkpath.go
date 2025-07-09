package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

func DoesPathExist(path string) (bool, error) {
	// Return the File/Path info
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// Creates All paths along the way
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Println("Error in MkdirAll: ", err)
			return false, err
		}
		// When there is an error
	} else if err != nil {
		fmt.Println("Error in deosPathExists: ", err)
		return false, err
	}
	// When the path exists
	return true, nil
}

func isDirEmpty(path string) (bool, error) {
	// Reads path into slice
	files, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	// If the slice is not empty
	if len(files) != 0 {
		return false, nil
	}

	// else need to be files
	return true, nil
}

func CheckOutputPath(filepath string) error {
	exists, err := DoesPathExist(filepath)
	if err != nil {
		return err
	}

	if exists {
		empty, err := isDirEmpty("output")
		if err != nil {
			return err
		}

		if !empty {
			if err := os.RemoveAll(filepath); err != nil {
				return err
			}
			if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkStringValidPath(path string) string {
	// Check if path is empty
	if len(path) == 0 {
		return ""
	}

	// Use filepath.Separator for cross-platform compatibility
	lastChar := path[len(path)-1:]
	if lastChar != string(filepath.Separator) {
		path += string(filepath.Separator)
	}

	_, err := DoesPathExist(path)
	if err != nil {
		return ""
	}
	return path
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

	// If no valid path was found, return an error
	if validPath == "" {
		return "", fmt.Errorf("no valid .mrpack file found for path: %s", path)
	}

	// Verify the file extension
	if filepath.Ext(validPath) != ".mrpack" {
		return "", fmt.Errorf("file does not have a .mrpack extension")
	}

	return validPath, nil
}
