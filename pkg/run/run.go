package run

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

const commandName = "run"

//go:embed *.tmpl*
var composeTemplates embed.FS

var Command *cli.Command = &cli.Command{
	Name:      commandName,
	Aliases:   []string{"r"},
	Usage:     "Run your application with Docker Compose.",
	UsageText: "3lv run [options] <application-name>",
	Flags: []cli.Flag{
		shared.SystemNameFlag(
			"The name of your system.",
		),
		shared.HelmValuesFileFlag(),
		shared.RegistryFlag(
			"The registry to use for the image. Used for finding the image name.",
		),
	},
	Action: Run,
}

func Run(_ context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	applicationName := c.Args().First()
	if applicationName == "" {
		return cli.Exit("Application name not provided", 1)
	}

	helmValues, err := parseHelmValuesFile(
		c.String("helm-values-file"),
	)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	composeFile, err := generateComposeFile(
		utils.StringWithDefault(c.String("registry"), "containerregistryelvia.azurecr.io"),
		c.String("system-name"),
		applicationName,
		helmValues,
	)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}

	dockerComposeUpOutput := dockerComposeUpCommand(
		composeFile,
		nil,
	)
	if command.IsError(dockerComposeUpOutput) {
		return cli.Exit(dockerComposeUpOutput.Error, 1)
	}

	return nil
}

type ComposeFileVariables struct {
	ApplicationName      string
	ImageName            string
	Port                 int
	TargetPort           int
	EnvironmentVariables map[string]string
}

func generateComposeFile(
	registry string,
	systemName string,
	applicationName string,
	helmValues *HelmValues,
) (string, error) {
	directory, err := os.MkdirTemp("", "3lv-run-*")
	if err != nil {
		return "", fmt.Errorf("Failed to create temporary directory: %w", err)
	}

	imageName, err := build.GetImageName(registry, systemName, applicationName)
	if err != nil {
		return "", err
	}

	composeFile, err := utils.WriteFileWithTemplate(
		directory,
		"docker-compose.yml",
		"docker-compose.yml.tmpl",
		composeTemplates,
		ComposeFileVariables{
			ApplicationName:      applicationName,
			ImageName:            imageName + ":latest-cache",
			Port:                 helmValues.Service.Port,
			TargetPort:           helmValues.Service.TargetPort,
			EnvironmentVariables: helmValues.GetEnvironmentVariablesMap(),
		},
	)
	if err != nil {
		return "", err
	}

	return composeFile, nil
}

func (v HelmValues) GetEnvironmentVariablesMap() map[string]string {
	env := make(map[string]string)
	for _, e := range v.Env {
		env[e.Name] = e.Value
	}

	return env
}

type HelmValues struct {
	Env []struct {
		Name  string `yaml:"name"`
		Value string `yaml:"value"`
	} `yaml:"env"`
	Service struct {
		Port       int `yaml:"port"`
		TargetPort int `yaml:"targetPort"`
	} `yaml:"service"`
}

func parseHelmValuesFile(
	helmValuesFilePath string,
) (*HelmValues, error) {
	var helmValues HelmValues
	if helmValuesFilePath == "" {
		return &helmValues, nil
	}

	helmValuesFileBytes, err := os.ReadFile(helmValuesFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(helmValuesFileBytes, &helmValues)
	if err != nil {
		return nil, err
	}

	return &helmValues, nil
}

func dockerComposeUpCommand(
	composeFile string,
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"docker",
			"compose",
			"-f",
			composeFile,
			"up",
		),
		options,
	)
}
