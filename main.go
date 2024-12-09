package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/create"
	"github.com/3lvia/cli/pkg/deploy"
	"github.com/3lvia/cli/pkg/githubactions"
	"github.com/3lvia/cli/pkg/run"
	"github.com/3lvia/cli/pkg/scan"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/upgrade"
	"github.com/urfave/cli/v3"
)

//go:embed VERSION
var version embed.FS

func main() {
	log.SetFlags(0)

	versionFile, err := version.ReadFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}
	version := strings.TrimSpace(string(versionFile))

	// Check for a new version of the CLI after all these commands.
	// Any new commands should be added here.
	commands := shared.WithCheckVersionAfterCommands(
		[]*cli.Command{
			build.Command,
			deploy.Command,
			scan.Command,
			githubactions.Command,
			run.Command,
			create.Command,
		},
	)

	app := &cli.Command{
		Name:                  "3lv",
		Usage:                 "Command Line Interface tool for developing, building and securing Elvia applications ⚡",
		Version:               version,
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "non-interactive",
				Usage:   "Run in non-interactive mode. Accepts all default values without prompting for confirmation.",
				Sources: cli.EnvVars("3LV_NON_INTERACTIVE"),
			},
		},
		// Don't check for updates when running the upgrade command.
		Commands: append(commands, upgrade.Command(version)),
	}

	ctx := context.Background()
	if err := app.Run(ctx, os.Args); err != nil {
		style.Print(
			fmt.Sprintf("\n\nERROR: %s", err),
			&style.PrintOptions{Color: "red"},
		)
	}
}
