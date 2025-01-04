package shared

import (
	"errors"
	"fmt"
	"os"

	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
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

const configFileName = "3lv.yml"

func GetConfig() (*Config, error) {
	exists, filePath := ConfigExists()
	if !exists {
		return &Config{},
			errors.New("No 3lv configuration file found, will proceed without it. You can run `3lv init` to create one.")
	}

	file, err := os.ReadFile(filePath)
	if err != nil {
		return &Config{},
			fmt.Errorf("Could not read 3lv configuration file at %s, will proceed without it.", filePath)
	}

	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return &Config{},
			fmt.Errorf("Could not parse 3lv configuration file at %s, will proceed without it.", filePath)
	}

	style.PrintInfo(fmt.Sprintf("Found 3lv confirguration file at %s.\n", filePath))

	return &config, nil
}

func SetConfig(config *Config) error {
	filePath := func() string {
		gitTopLevel, err := utils.ResolveGitRepositoryTopLevelPath()
		if err != nil {
			return configFileName
		}

		return gitTopLevel + "/" + configFileName
	}()

	file, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, file, 0o644)

	return err
}

func ConfigExists() (bool, string) {
	filePath := func() string {
		gitTopLevel, err := utils.ResolveGitRepositoryTopLevelPath()
		if err != nil {
			return configFileName
		}

		return gitTopLevel + "/" + configFileName
	}()

	_, err := os.Stat(filePath)

	return !os.IsNotExist(err), filePath
}

func WithConfig(commands [](func(_ *Config) *cli.Command), config *Config) []*cli.Command {
	return lo.Map(commands, func(c func(_ *Config) *cli.Command, _ int) *cli.Command {
		return c(config)
	})
}

func (config *Config) GetConfigForApplication(applicationName string) (*Application, error) {
	if config == nil {
		return &Application{}, errors.New("Config is nil.")
	}

	for _, application := range config.Applications {
		if application.Name == applicationName {
			return &application, nil
		}
	}

	return &Application{}, fmt.Errorf("Application %s not found in config.", applicationName)
}

func (config *Config) IsEmpty() bool {
	return config.System == "" && len(config.Applications) == 0
}
