package utils

import (
	"os"
)

func CreateFile(path string) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		if IsFileExists(path) {
			err = os.Remove(path)
			if err != nil {
				return nil, err
			}
			return CreateFile(path)
		}

		return nil, err
	}

	err = os.Chmod(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func IsFileExists(path string) bool {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}
