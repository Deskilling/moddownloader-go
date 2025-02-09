package main

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
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
		return allFiles, nil
	}
	return nil, nil
}

func calculateHashSha1(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func calculateHashSha512(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func calcualteAllHashesFromDirectory(directory string) ([]string, []string, error) {
	allFiles, err := getAllFilesFromDirectory(directory)
	if err != nil {
		return nil, nil, err
	}

	var sha1Hashes []string
	var sha512Hashes []string
	for _, file := range allFiles {
		filepath := directory + file.Name()
		//hash512, err := calculateHashSha512(filepath)
		hash1, err := calculateHashSha1(filepath)
		if err != nil {
			return nil, nil, err
		}
		sha1Hashes = append(sha1Hashes, hash1)

		hash512, err := calculateHashSha512(filepath)
		if err != nil {
			return nil, nil, err
		}
		sha512Hashes = append(sha512Hashes, hash512)

	}

	return sha1Hashes, sha512Hashes, nil
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
			os.RemoveAll(filepath)
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	return nil
}
