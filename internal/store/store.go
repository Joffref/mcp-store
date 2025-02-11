package store

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

type Store struct {
	Repositories map[string]*Repository `yaml:"repositories"`
}

type Repository struct {
	Repository      string   `yaml:"repository" mendatory:"true"`
	SmitheryPath    string   `yaml:"smitheryPath" mendatory:"false" default:"smithery.yaml"`
	Dockerfile      string   `yaml:"dockerfile" mendatory:"false" default:"Dockerfile"`
	Branch          string   `yaml:"branch" mendatory:"false" default:"main"`
	DisplayName     string   `yaml:"displayName" mendatory:"true"`
	Icon            string   `yaml:"icon" mendatory:"true"`
	Description     string   `yaml:"description" mendatory:"true"`
	LongDescription string   `yaml:"longDescription" mendatory:"true"`
	Tags            []string `yaml:"tags"`
	Categories      []string `yaml:"categories"`
}

func (s *Store) Read(path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlFile, s)
}

// ValidateWithDefaultValues validates the store and applies default values to empty fields
// This is useful to validate the store before running the import command
func (s *Store) ValidateWithDefaultValues() error {
	if s.Repositories == nil {
		return errors.New("repositories is required")
	}

	var errs []error

	for name, repository := range s.Repositories {
		// Use reflection to validate struct tags
		v := reflect.ValueOf(repository).Elem() // Get the element the pointer refers to
		t := v.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			// Check mandatory fields
			if mandatory, ok := field.Tag.Lookup("mendatory"); ok && mandatory == "true" {
				if value.IsZero() {
					errs = append(errs, fmt.Errorf("field %s is required in repository %s", field.Name, name))
				}
			}

			// Apply default values for empty fields
			if defaultVal, ok := field.Tag.Lookup("default"); ok && value.IsZero() {
				switch value.Kind() {
				case reflect.String:
					value.SetString(defaultVal)
				}
			}
		}
	}

	return errors.Join(errs...)
}
