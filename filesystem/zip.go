package filesystem

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ExtractZip(source, dest string) error {
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

// Kuss gosamples
func ZipSource(source, target string) error {
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
