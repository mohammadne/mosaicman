package utils

import "os"

func CreateDirIfMissed(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return os.MkdirAll(dirName, os.ModePerm)
	}

	return nil
}
