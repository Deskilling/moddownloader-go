package filesystem

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func ExtractZip(source, dest string) error {
	read, err := zip.OpenReader(source)
	if err != nil {
		log.Error("failed reading zip", "source", source, "err", err)
		return err
	}
	defer read.Close()
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			log.Error("failed opening file", "file", file, "err", err)
			return err
		}
		name := filepath.Join(dest, file.Name)
		err = CreatePath(filepath.Dir(name))
		if err != nil {
			log.Error("failed creating path", "path", filepath.Dir(name), "err", err)
			return err
		}
		create, err := os.Create(name)
		if err != nil {
			log.Error("failed creating file", "file", name, "err", err)
			return err
		}
		defer create.Close()
		_, err = io.Copy(create, open)
		if err != nil {
			log.Error("failed copying file", "file", create, "err", err)
			return err
		}
		open.Close()
	}
	return nil
}

func ZipSource(source, target string) error {
	if _, err := os.Stat(source); err != nil {
		log.Error("source does not exist", "err", err)
		return err
	}

	f, err := os.Create(target)
	if err != nil {
		log.Error("failed to create zip file", "err", err)
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("walk error", "path", path, "err", err)
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			log.Error("failed to create info header", "err", err)
			return err
		}

		header.Method = zip.Deflate

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			log.Error("failed to calculate relative path", "source", source, "path", path, "err", err)
			return err
		}
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			log.Error("failed to create header in zip", "header", header, "err", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			log.Error("ailed to open file", "path", path, "err", err)
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		if err != nil {
			log.Error("failed to write file contents", "err", err)
			return err
		}

		return nil
	})
}
