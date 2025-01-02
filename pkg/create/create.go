package create

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strings"

	"github.com/3lvia/cli/pkg/build"
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
	Dotnet8WebAPI = Template{"dotnet8-webapi"}
	// Dotnet8WebApp = Template{"dotnet8-webapp"}.
	Dotnet8Worker = Template{"dotnet8-worker"}
	// Go           Template = Template{"go"}.
	PythonWebAPI = Template{"python-webapi"}
	Templates    = enum.New(
		Dotnet8WebAPI,
		//	Dotnet8WebApp,
		Dotnet8Worker,
		// Go,
		PythonWebAPI,
	)
)

var Command *cli.Command = &cli.Command{
	Name:      commandName,
	Aliases:   []string{"c"},
	Usage:     "Create a new project from one of Elvia's templates.",
	UsageText: "3lv create [options] <output-directory>",
	Flags: []cli.Flag{
		shared.SystemNameFlag(
			"The name of your system (Kubernetes namespace) you want to create your application in.",
		),
		shared.ApplicationNameFlag(
			"The name of the application you want to create.",
		),
		shared.RuntimeCloudProviderFlag(),
		&cli.StringFlag{
			Name:    "template",
			Aliases: []string{"t"},
			Usage:   "The template to use for the project. Supported templates are: " + strings.Join(Templates.Values(), ", "),
			Value:   Dotnet8WebAPI.Value,
			Action: func(_ context.Context, _ *cli.Command, template string) error {
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
			Usage: "The root directory of your GitHub repository." +
				" The path specified will be prepended to '.github/workflows'.",
		},
		&cli.StringFlag{
			Name:  "python-version",
			Usage: "The version of Python to use for the project. Only applicable for Python templates.",
		},
	},
	Action: Create,
}

func Create(ctx context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

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
	nonInteractive := c.Bool("non-interactive")
	pythonVersion := c.String("python-version")

	if template != PythonWebAPI && c.IsSet("python-version") {
		style.PrintWarning("Argument 'python-version' is only applicable for Python templates.")
	}

	checkCoooiecutterInstalledOutput := checkCookiecutterInstalledCommand(nil)
	if command.IsError(checkCoooiecutterInstalledOutput) {
		yes, err := utils.PromptYesNo("Cookiecutter is not installed. Do you want to install it?", nonInteractive)
		if err != nil {
			return cli.Exit(err, 1)
		}

		if yes {
			checkPipxInstalledOutput := checkPipxInstalledCommand(nil)
			if command.IsError(checkPipxInstalledOutput) {
				log.Fatal("pipx, which is required for installing cookiecutter, is not installed. Please install it first.")
			}

			style.PrintInfo("Installing cookiecutter...")

			installCookiecutterOutput := installCookiecutterCommand(nil)
			if command.IsError(installCookiecutterOutput) {
				return cli.Exit("Failed to install cookiecutter.", 1)
			}

			style.PrintSuccess("Cookiecutter installed!")
		} else {
			return cli.Exit("Cookiecutter is required for creating a new project. Please install it first.", 1)
		}
	}

	cookiecutterOutput := cookiecutterCommand(
		template,
		outputDirectory,
		applicationName,
		systemName,
		pythonVersion,
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

	if template == PythonWebAPI {
		checkUvInstalledOutput := checkUvInstalledCommand(nil)
		if command.IsError(checkUvInstalledOutput) {
			yes, err := utils.PromptYesNo("uv is not installed. Do you want to install it using pipx?", nonInteractive)
			if err != nil {
				return cli.Exit(err, 1)
			}

			checkPipxInstalledOutput := checkPipxInstalledCommand(nil)
			if command.IsError(checkPipxInstalledOutput) {
				log.Fatal(
					"pipx is not installed, cannot automatically install uv." +
						" Please install uv yourself (https://docs.astral.sh/uv/getting-started/installation)," +
						" or install pipx and try again.",
				)
			}

			if yes {
				style.PrintInfo("Installing uv...")

				installUvOutput := installUvCommand(nil)
				if command.IsError(installUvOutput) {
					return cli.Exit("Failed to install uv.", 1)
				}

				style.PrintSuccess("uv installed!")
			} else {
				return cli.Exit("uv is required for creating a new project. Please install it first.", 1)
			}
		}

		uvSyncOutput := uvSyncCommand(projectDirectory, nil)
		if command.IsError(uvSyncOutput) {
			return cli.Exit("Failed to generate uv.lock file.", 1)
		}
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
		ctx,
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

	style.PrintSuccess(fmt.Sprintf("Successfully created project at '%s'!", projectDirectory))

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
	case Dotnet8WebAPI /*Dotnet8WebApp,*/, Dotnet8Worker:
		return path.Join(
			outputDirectory,
			toPascalCaseWithoutHyphens(applicationName),
		), nil
	case PythonWebAPI /*, Go*/ :
		return path.Join(outputDirectory, applicationName), nil
	default:
		return "", fmt.Errorf("Could not find project directory for template '%s'", template)
	}
}

func getProjectFileForTemplate(
	template Template,
	applicationName string,
) (string, error) {
	switch template {
	case Dotnet8WebAPI /*Dotnet8WebApp,*/, Dotnet8Worker:
		return toPascalCaseWithoutHyphens(applicationName) + ".csproj", nil
	/*
		case Go:
			return "go.mod", nil
	*/
	case PythonWebAPI:
		return "pyproject.toml", nil
	default:
		return "", fmt.Errorf("Could not find project file for template '%s'", template)
	}
}

func cookiecutterCommand(
	template Template,
	outputDirectory string,
	applicationName string,
	systemName string,
	pythonVersion string,
	options *command.RunOptions,
) command.Output {
	cmd := *exec.Command(
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
	)

	if template == PythonWebAPI {
		if pythonVersion == "" {
			cmd.Args = append(cmd.Args, "python_version="+build.DefaultPythonVersion)
		} else {
			cmd.Args = append(cmd.Args, "python_version="+pythonVersion)
		}
	}

	return command.Run(cmd, options)
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

func checkPipxInstalledCommand(
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

func checkUvInstalledCommand(
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"uv",
			"--version",
		),
		options,
	)
}

func installUvCommand(
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"pipx",
			"install",
			"uv",
		),
		options,
	)
}

func uvSyncCommand(
	projectDirectory string,
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"uv",
			"sync",
			"--directory",
			projectDirectory,
		),
		options,
	)
}
