package create

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/githubactions"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/orsinium-labs/enum"
	"github.com/urfave/cli/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const commandName = "create"

type Template enum.Member[string]

var (
	Dotnet8WebApi Template = Template{"dotnet8-webapi"}
	// Dotnet8WebApp Template = Template{"dotnet8-webapp"}
	Dotnet8Worker Template = Template{"dotnet8-worker"}
	// Go            Template = Template{"go"}
	Templates = enum.New(
		Dotnet8WebApi,
		//	Dotnet8WebApp,
		Dotnet8Worker,
		// Go,
	)
)

var Command *cli.Command = &cli.Command{
	Name:    commandName,
	Aliases: []string{"c"},
	Usage:   "Create a new project",
	Flags: []cli.Flag{
		shared.SystemNameFlag(
			"The name of your system (Kubernetes namespace) you want to create your application in.",
			true,
		),
		shared.ApplicationNameFlag(
			"The name of the application you want to create.",
		),
		shared.RuntimeCloudProviderFlag(),
		&cli.StringFlag{
			Name:    "template",
			Aliases: []string{"t"},
			Usage: fmt.Sprintf(
				"The template to use for the project. Supported templates are: %s",
				strings.Join(Templates.Values(), ", "),
			),
			Value: Dotnet8WebApi.Value,
			Action: func(ctx context.Context, c *cli.Command, template string) error {
				parsed := Templates.Parse(template)

				if parsed == nil {
					return cli.Exit(
						fmt.Sprintf(
							"Template '%s' is not supported. Supported templates are: %s",
							template,
							Templates.Values(),
						),
						1,
					)
				}

				return nil
			},
		},
		&cli.StringFlag{
			Name:    "default-branch",
			Aliases: []string{"b"},
			Usage:   "The default branch of the repository",
			Value:   "trunk",
		},
		&cli.StringFlag{
			Name:    "github-actions-directory",
			Aliases: []string{"G"},
			Usage:   "The root directory of your GitHub repository. The path specified will be prepended to '.github/workflows'.",
		},
	},
	Action: Create,
}

func Create(ctx context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		return cli.ShowAppHelp(c)
	}

	nonInteractive := c.Bool("non-interactive")

	outputDirectory := c.Args().First()
	if outputDirectory == "" {
		return cli.Exit("Output directory not provided", 1)
	}

	systemName := c.String("system-name")
	if systemName == "" {
		return cli.Exit("System name not provided", 1)
	}

	applicationName := c.String("application-name")
	if applicationName == "" {
		return cli.Exit("Application name not provided", 1)
	}

	template := func() Template {
		parsed := Templates.Parse(c.String("template"))
		if parsed == nil {
			return Template{"dotnet"}
		}
		return *parsed
	}()
	defaultBranch := c.String("default-branch")

	checkCoooiecutterInstalledOutput := checkCookiecutterInstalledCommand(nil)
	if command.IsError(checkCoooiecutterInstalledOutput) {
		yes, err := utils.PromptYesNo("Cookiecutter is not installed. Do you want to install it?", nonInteractive)
		if err != nil {
			return cli.Exit(err, 1)
		}

		if yes {
			checkPipxInstalledOutput := checkPipxÌnstalledCommand(nil)
			if command.IsError(checkPipxInstalledOutput) {
				log.Fatal("pipx, which is required for installing cookiecutter, is not installed. Please install it first.")
			}

			style.Print(
				"Installing cookiecutter...",
				nil,
			)
			installCookiecutterOutput := installCookiecutterCommand(nil)
			if command.IsError(installCookiecutterOutput) {
				return cli.Exit("Failed to install cookiecutter.", 1)
			}

			style.Print(
				"Cookiecutter installed!",
				&style.PrintOptions{Color: "green"},
			)
		} else {
			return cli.Exit("Cookiecutter is required for creating a new project. Please install it first.", 1)
		}
	}

	cookiecutterOutput := cookiecutterCommand(
		template,
		outputDirectory,
		applicationName,
		systemName,
		nil,
	)

	if command.IsError(cookiecutterOutput) {
		return cli.Exit("Failed to create project.", 1)
	}

	projectDirectory, err := getProjectDirectoryForTemplate(
		template,
		outputDirectory,
		applicationName,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	githubActionsDirectory := func() string {
		if c.IsSet("github-actions-directory") {
			return c.String("github-actions-directory")
		}
		return projectDirectory
	}()

	projectFile, err := getProjectFileForTemplate(template, applicationName)
	if err != nil {
		return cli.Exit(err, 1)
	}

	err = githubactions.CreateDeployWorkflow(
		githubActionsDirectory,
		projectFile,
		c.String("runtime-cloud-provider"),
		systemName,
		applicationName,
		"", // intentional
		defaultBranch,
		nonInteractive,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	style.Print(
		fmt.Sprintf("Succesfully created project at '%s'!", projectDirectory),
		&style.PrintOptions{Color: "green"},
	)

	return nil
}

func toPascalCaseWithoutHyphens(s string) string {
	return strings.ReplaceAll(cases.Title(language.English).String(s), "-", "")
}

func getProjectDirectoryForTemplate(
	template Template,
	outputDirectory string,
	applicationName string,
) (string, error) {
	switch template {
	case Dotnet8WebApi /*Dotnet8WebApp,*/, Dotnet8Worker:
		return path.Join(
			outputDirectory,
			toPascalCaseWithoutHyphens(applicationName),
		), nil
	/*
		case Go:
			return path.Join(outputDirectory, applicationName), nil
	*/
	default:
		return "", fmt.Errorf("Could not find project directory for template '%s'", template)
	}
}

func getProjectFileForTemplate(
	template Template,
	applicationName string,
) (string, error) {
	switch template {
	case Dotnet8WebApi /*Dotnet8WebApp,*/, Dotnet8Worker:
		return fmt.Sprintf(
			"%s.csproj",
			toPascalCaseWithoutHyphens(applicationName),
		), nil
	/*
		case Go:
			return "go.mod", nil
	*/
	default:
		return "", fmt.Errorf("Could not find project file for template '%s'", template)
	}
}

func cookiecutterCommand(
	template Template,
	outputDirectory string,
	applicationName string,
	systemName string,
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"cookiecutter",
			"gh:3lvia/application-templates",
			"--directory",
			template.Value,
			"--output-dir",
			outputDirectory,
			"--no-input",
			"application_name="+applicationName,
			"application_name_pascal_case="+toPascalCaseWithoutHyphens(applicationName),
			"system_name="+systemName,
			// TODO: is this needed?
			"base_dir=./",
		),
		options,
	)
}

func checkCookiecutterInstalledCommand(
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"cookiecutter",
			"--version",
		),
		options,
	)
}

func checkPipxÌnstalledCommand(
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"pipx",
			"--version",
		),
		options,
	)
}

func installCookiecutterCommand(
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"sudo",
			"pipx",
			"install",
			"cookiecutter",
			"--global",
		),
		options,
	)
}
