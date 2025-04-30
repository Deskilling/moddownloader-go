package main

import (
	"archive/zip"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func doesPathExist(path string) (bool, error) {
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

func getAllFilesFromDirectory(directory string, extension string) ([]os.DirEntry, error) {
	doesExist, err := doesPathExist(directory)
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

func calculateHashes(filepath string) (string, string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	sha1Hash := sha1.New()
	sha512Hash := sha512.New()

	// Copy file content to both hash functions simultaneously
	if _, err := io.Copy(io.MultiWriter(sha1Hash, sha512Hash), file); err != nil {
		return "", "", err
	}

	return hex.EncodeToString(sha1Hash.Sum(nil)), hex.EncodeToString(sha512Hash.Sum(nil)), nil
}

func calculateAllHashesFromDirectory(directory string) ([]string, []string, []os.DirEntry, error) {
	allFiles, err := getAllFilesFromDirectory(directory, ".jar")
	if err != nil {
		return nil, nil, nil, err
	}

	var sha1Hashes, sha512Hashes []string

	for _, file := range allFiles {
		filePath := filepath.Join(directory, file.Name())

		hash1, hash512, err := calculateHashes(filePath)
		if err != nil {
			return nil, nil, nil, err
		}

		sha1Hashes = append(sha1Hashes, hash1)
		sha512Hashes = append(sha512Hashes, hash512)
	}

	return sha1Hashes, sha512Hashes, allFiles, nil
}

func checkOutputPath(filepath string) error {
	exists, err := doesPathExist(filepath)
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

func extractZip(source, dest string) error {
	read, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer read.Close()
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := filepath.Join(dest, file.Name)
		err = os.MkdirAll(filepath.Dir(name), os.ModeDir)
		if err != nil {
			return err
		}
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		_, err = io.Copy(create, open)
		if err != nil {
			return err
		}
		open.Close()
	}
	return nil
}

func readFile(filepath string) string {
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}
	return string(fileContent)
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

	_, err := doesPathExist(path)
	if err != nil {
		return ""
	}
	return path
}

// checkMrpack verifies if the given path points to a valid .mrpack file
func checkMrpack(path string) (string, error) {
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

func writeFile(path string, content []byte) {
	// 0064 is the permission level
	err := os.WriteFile(path, content, 0064)
	if err != nil {
		return
	}
}

// Kuss gosamples
func zipSource(source, target string) error {
	// Ensure source path exists
	if _, err := os.Stat(source); err != nil {
		return fmt.Errorf("source path error: %w", err)
	}

	// Create target file
	f, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// Walk through source directory
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create header: %w", err)
		}

		// Set compression method
		header.Method = zip.Deflate

		// Calculate relative path
		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return fmt.Errorf("failed to calculate relative path: %w", err)
		}
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store // Don't compress directories
		}

		// Create header in zip
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create header in zip: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		// Open and copy file contents
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		if err != nil {
			return fmt.Errorf("failed to write file contents: %w", err)
		}

		return nil
	})
}
