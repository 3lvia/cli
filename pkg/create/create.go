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
	"github.com/urfave/cli/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const commandName = "create"

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
			Usage:   "The template to use for the project",
			Value:   "dotnet",
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

	templateName := c.String("template")
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

	applicationNamePascalCase := strings.ReplaceAll(cases.Title(language.English).String(applicationName), "-", "")

	cookiecutterOutput := cookiecutterCommand(
		templateName,
		outputDirectory,
		applicationName,
		applicationNamePascalCase,
		systemName,
		nil,
	)

	if command.IsError(cookiecutterOutput) {
		return cli.Exit("Failed to create project.", 1)
	}

	// TODO: depend on template
	projectDirectory := path.Join(outputDirectory, applicationNamePascalCase)

	githubActionsDirectory := func() string {
		if c.IsSet("github-actions-directory") {
			return c.String("github-actions-directory")
		}
		return projectDirectory
	}()

	err := githubactions.CreateDeployWorkflow(
		githubActionsDirectory,
		// TODO: depend on template
		path.Join(projectDirectory, applicationNamePascalCase+".csproj"),
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

func cookiecutterCommand(
	templateName string,
	outputDirectory string,
	applicationName string,
	applicationNamePascalCase string,
	systemName string,
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"cookiecutter",
			"gh:3lvia/application-templates",
			"--directory",
			templateName,
			"--output-dir",
			outputDirectory,
			"--no-input",
			"application_name="+applicationName,
			"application_name_pascal_case="+applicationNamePascalCase,
			"system_name="+systemName,
			"bff_client_name="+applicationName,
			"sub_domain_name="+applicationName,
			"domain_path=/api",
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
			"pipx",
			"install",
			"cookiecutter",
		),
		options,
	)
}
