package shared

import (
	"errors"
	"fmt"
	"os"
	"path"

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
	Name                   string `yaml:"name"`
	ProjectFile            string `yaml:"projectFile"`
	BuildContext           string `yaml:"buildContext"`
	HelmValuesFile         string `yaml:"helmValuesFile"`
	GoMainPackageDirectory string `yaml:"goMainPackageDirectory,omitempty"`
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

func SetConfig(config *Config, overrideDirectory string) error {
	filePath := func() string {
		if overrideDirectory != "" {
			return path.Join(overrideDirectory, configFileName)
		}

		gitTopLevel, err := utils.ResolveGitRepositoryTopLevelPath()
		if err != nil {
			return configFileName
		}

		return path.Join(gitTopLevel, configFileName)
	}()

	file, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, file, 0o644)

	return err
}

func ConfigExists() (bool, string) {
	returnNameCheckExists := func(filePath string) (string, error) {
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			return "", errors.New("No 3lv configuration file found at " + filePath)
		}

		return filePath, nil
	}

	gitFilePath, gitErr := func() (string, error) {
		gitTopLevel, err := utils.ResolveGitRepositoryTopLevelPath()
		if err != nil {
			return "", err
		}

		return returnNameCheckExists(path.Join(gitTopLevel, configFileName))
	}()

	currentDirectoryFilePath, currentDirectoryErr := func() (string, error) {
		currentDirectory, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return returnNameCheckExists(path.Join(currentDirectory, configFileName))
	}()

	defaultFilePath, defaultErr := returnNameCheckExists(configFileName)

	if gitErr == nil {
		return true, gitFilePath
	}

	if currentDirectoryErr == nil {
		return true, currentDirectoryFilePath
	}

	if defaultErr == nil {
		return true, defaultFilePath
	}

	return false, ""
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
