package config

import (
	"os"
	"path"

	"github.com/mlange-42/tom/api"
	"github.com/mlange-42/tom/util"
	"gopkg.in/yaml.v3"
)

func LoadLocations() (map[string]api.Location, error) {
	dir, err := GetRootDir()
	if err != nil {
		return nil, err
	}

	path := path.Join(dir, locationsFile)

	if !util.FileExists(path) {
		return map[string]api.Location{}, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	locs := map[string]api.Location{}

	if err := yaml.Unmarshal(file, &locs); err != nil {
		return nil, err
	}

	return locs, nil
}

func SaveLocations(locations map[string]api.Location) error {
	dir, err := GetRootDir()
	if err != nil {
		return err
	}
	err = CreateDir(dir)
	if err != nil {
		return err
	}

	path := path.Join(dir, locationsFile)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := yaml.Marshal(&locations)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}
