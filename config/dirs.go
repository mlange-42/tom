package config

import (
	"os"
	"path/filepath"
)

const appDirectory = ".tom"
const locationsFile = "locations.yml"

func GetRootDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, appDirectory), nil
}

func CreateDir(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
