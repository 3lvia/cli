package shared

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

type Config struct {
	System       string        `yaml:"system"`
	Applications []Application `yaml:"applications"`
}

type Application struct {
	Name           string `yaml:"name"`
	BuildContext   string `yaml:"buildContext,omitempty"`
	ProjectFile    string `yaml:"projectFile"`
	HelmValuesFile string `yaml:"helmValuesFile"`
}

func ReadConfig() (*Config, error) {
	// TODO: read from git root
	file, err := os.ReadFile("3lv.yml")
	if err != nil {
		return &Config{}, err
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return &Config{}, err
	}

	return &config, nil
}

func WithConfig(commands [](func(_ *Config) *cli.Command), config *Config) []*cli.Command {
	return lo.Map(commands, func(c func(_ *Config) *cli.Command, _ int) *cli.Command {
		return c(config)
	})
}

func (config Config) GetConfigForApplication(applicationName string) (*Application, error) {
	for _, application := range config.Applications {
		if application.Name == applicationName {
			return &application, nil
		}
	}

	return &Application{}, fmt.Errorf("application %s not found in config", applicationName)
}
