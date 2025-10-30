package filesystem

import (
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func CalculateHashes(filepath string) (string, string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	sha1Hash := sha1.New()
	sha512Hash := sha512.New()

	if _, err := io.Copy(io.MultiWriter(sha1Hash, sha512Hash), file); err != nil {
		return "", "", err
	}

	return hex.EncodeToString(sha1Hash.Sum(nil)), hex.EncodeToString(sha512Hash.Sum(nil)), nil
}

func CalculateDirectoryHashes(directory, extension string) ([]string, []string, []os.DirEntry, error) {
	allFiles, err := ReadDirectory(directory, extension)
	if err != nil {
		log.Error("failed reading", "directory", directory, "err", err)
		return nil, nil, nil, err
	}

	var sha1Hashes, sha512Hashes []string
	for _, file := range allFiles {
		filePath := filepath.Join(directory, file.Name())

		hash1, hash512, err := CalculateHashes(filePath)
		if err != nil {
			log.Error("failed calculating hash", "filepath", filePath, "err", err)
			return nil, nil, nil, err
		}

		sha1Hashes = append(sha1Hashes, hash1)
		sha512Hashes = append(sha512Hashes, hash512)
	}

	return sha1Hashes, sha512Hashes, allFiles, nil
}
