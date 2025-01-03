package initialize

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/3lvia/cli/pkg/shared"
	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

const commandName = "init"

func Command() *cli.Command {
	return &cli.Command{
		Name:      commandName,
		Aliases:   []string{"i"},
		Usage:     "Create a new 3lv configuration file.",
		UsageText: "3lv init",
		Action:    Init,
	}
}

func Init(_ context.Context, c *cli.Command) error {
	if c.IsSet("non-interactive") {
		return cli.Exit("Non-interactive mode is not supported for this command.", 1)
	}

	var systemName string

	err := huh.NewInput().
		Title("What is the name of your system? (e.g. 'core')").
		Value(&systemName).
		Run()
	if err != nil {
		return err
	}

	var applications []shared.Application

	for {
		var applicationName string

		err := huh.NewInput().
			Title("What is the name of your application? (e.g. 'demo-api')").
			Value(&applicationName).
			Run()
		if err != nil {
			return err
		}

		var projectFile string

		err = huh.NewInput().
			Title("What is the full path to your project file? (e.g. 'applications/demo-api/demo-api.csproj')").
			Description("Supported project files are *.csproj, go.mod, pyproject.toml, and Dockerfile").
			Value(&projectFile).
			Validate(func(str string) error {
				base := path.Base(str)

				if !(strings.HasSuffix(base, ".csproj") ||
					base == "go.mod" ||
					base == "pyproject.toml" ||
					strings.Contains(base, "Dockerfile")) {
					return fmt.Errorf("'%s' is not a supported project file", base)
				}

				return nil
			}).
			Run()
		if err != nil {
			return err
		}

		helmValuesFile := ".github/deploy/values-" + applicationName + ".yml"

		err = huh.NewInput().
			Title("What is the full path to your Helm values file?").
			Value(&helmValuesFile).
			Run()
		if err != nil {
			return err
		}

		var anotherOne bool

		err = huh.NewConfirm().
			Title("Do you want to add another application?").
			Value(&anotherOne).
			Run()
		if err != nil {
			return err
		}

		applications = append(applications, shared.Application{ //nolint:exhaustruct
			Name:           applicationName,
			ProjectFile:    projectFile,
			HelmValuesFile: helmValuesFile,
		})

		if !anotherOne {
			break
		}
	}

	config := &shared.Config{
		System:       systemName,
		Applications: applications,
	}

	exists, _ := shared.ConfigExists()
	if exists {
		var overwrite bool

		err = huh.NewConfirm().
			Title("A 3lv configuration file already exists, do you want to overwrite it?").
			Value(&overwrite).
			Run()
		if err != nil {
			return err
		}

		if !overwrite {
			return nil
		}
	}

	err = shared.SetConfig(config)

	return err
}
