package main

import (
	"archive/zip"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
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

func getAllFilesFromDirectory(directory string) ([]os.DirEntry, error) {
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
			if filepath.Ext(file.Name()) == ".jar" {
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
	allFiles, err := getAllFilesFromDirectory(directory)
	if err != nil {
		return nil, nil, nil, err
	}

	var sha1Hashes, sha512Hashes []string

	for _, file := range allFiles {
		filePath := directory + file.Name()

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
		name := path.Join(dest, file.Name)
		err = os.MkdirAll(path.Dir(name), os.ModeDir)
		if err != nil {
			return err
		}
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
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
	lastChar := path[len(path)-1:]
	if lastChar != "/" {
		path += "/"
	}
	doesPathExist(path)
	return path
}
