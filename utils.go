package gin_docs

import (
	"io"
	"os"
	"path/filepath"
)

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func copyFolder(source, destination string) error {
	files, err := filepath.Glob(filepath.Join(source, "*"))
	if err != nil {
		return err
	}

	err = os.MkdirAll(destination, os.ModePerm)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return err
		}

		destinationFile := filepath.Join(destination, fileInfo.Name())

		if fileInfo.IsDir() {
			err = copyFolder(file, destinationFile)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(file, destinationFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}
