package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

func zipit(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()


	info, err := os.Stat(source)
	if err != nil {
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
		fmt.Println("baseDir:"+baseDir)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		fmt.Println("path:"+path)

		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join("", strings.TrimPrefix(path, source))
			fmt.Println("header.Name:" + header.Name)
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
func main() {
	zipit("./documents", "./backup.zip")
}
