package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Repositories map[string]Repository `yaml:"repositories"`
}

type Repository struct {
	Repository   string `yaml:"repository"`
	SmitheryPath string `yaml:"smitheryPath"`
	Branch       string `yaml:"branch"`
	DisplayName  string `yaml:"displayName"`
	Icon         string `yaml:"icon"`
	Description  string `yaml:"description"`
}

func (c *Config) Read(path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, c)
}

func (c *Config) Validate() error {
	if c.Repositories == nil {
		return errors.New("repositories is required")
	}

	for _, repository := range c.Repositories {
		if repository.Repository == "" {
			return errors.New("repository is required")
		}

		if repository.DisplayName == "" {
			return errors.New("displayName is required")
		}

		if repository.Icon == "" {
			return errors.New("icon is required")
		}

		if repository.Description == "" {
			return errors.New("description is required")
		}
	}

	return nil
}
