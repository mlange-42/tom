package config

import (
	"os"
	"path"

	"github.com/mlange-42/tom/util"
	"gopkg.in/yaml.v3"
)

type CliArgs struct {
	Location   string
	Coords     Location `yaml:"-"`
	SetDefault bool     `yaml:"-"`
	Days       int
	PastDays   int
	Service    Service
}

func LoadCliArgs() (CliArgs, error) {
	dir, err := GetRootDir()
	if err != nil {
		return CliArgs{}, err
	}

	path := path.Join(dir, configFile)

	if !util.FileExists(path) {
		return CliArgs{
			Days:    7,
			Service: Services[0],
		}, nil
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return CliArgs{}, err
	}

	args := CliArgs{}

	if err := yaml.Unmarshal(file, &args); err != nil {
		return CliArgs{}, err
	}

	return args, nil
}

func SaveCliArgs(args *CliArgs) error {
	dir, err := GetRootDir()
	if err != nil {
		return err
	}
	err = CreateDir(dir)
	if err != nil {
		return err
	}

	path := path.Join(dir, configFile)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := yaml.Marshal(args)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}
