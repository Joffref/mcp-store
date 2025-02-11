package smithery

import (
	"fmt"
	"strings"
)

type SmitheryConfig struct {
	ParsedCommand *Command     `yaml:"parsedConfig,omitempty"`
	Build         *Build       `yaml:"build,omitempty"`
	StartCommand  StartCommand `yaml:"startCommand"`
}

type Build struct {
	Dockerfile      *string `yaml:"dockerfile,omitempty"`
	DockerBuildPath *string `yaml:"dockerBuildPath,omitempty"`
}

type Command struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

func (c *Command) String() string {
	return fmt.Sprintf("%s %s", c.Command, strings.Join(c.Args, " "))
}

type StartCommand struct {
	Type            string       `yaml:"type"`
	ConfigSchema    ConfigSchema `yaml:"configSchema"`
	CommandFunction string       `yaml:"commandFunction"`
}

type ConfigSchema struct {
	Type       string              `yaml:"type"`
	Required   []string            `yaml:"required"`
	Properties map[string]Property `yaml:"properties"`
}

type Property struct {
	Type        string `yaml:"type"`
	Default     string `yaml:"default"`
	Description string `yaml:"description"`
}
