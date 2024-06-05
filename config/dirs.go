package config

import (
	"os"
)

const locationsFile = "tom/locations.yml"
const configFile = "tom/config.yml"

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
