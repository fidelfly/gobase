package filex

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// This function is to make sure that the folder exist.
func MustFoler(path string) error {
	if _, err := os.Stat(path); err == nil {
		//if info.Mode().Perm()
		//os.Chmod(path, os.ModePerm)
	} else {
		err := os.MkdirAll(path, os.ModeDir)
		if err != nil {
			return err
		}
		_ = os.Chmod(path, os.ModePerm)

	}
	return nil
}

func ZipFile(source, target string) error {
	zf, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zf.Close()

	archive := zip.NewWriter(zf)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	_ = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}
