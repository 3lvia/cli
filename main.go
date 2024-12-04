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
	"github.com/3lvia/cli/pkg/style"
	"github.com/google/go-github/v66/github"
	"github.com/urfave/cli/v3"
	"golang.org/x/mod/semver"
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
		After: func(ctx context.Context, c *cli.Command) error {
			latestVersion, _ := getLatestVersion(ctx)
			if semver.Compare("v"+version, "v"+latestVersion) == -1 {
				style.Print(
					fmt.Sprintf("\n\nA new version of 3lv is available! %s -> %s", version, latestVersion),
					&style.PrintOptions{Color: "yellow"},
				)
			}

			return nil
		},
		Commands: []*cli.Command{
			build.Command,
			deploy.Command,
			scan.Command,
			githubactions.Command,
			run.Command,
			create.Command,
		},
	}

	ctx := context.Background()
	if err := app.Run(ctx, os.Args); err != nil {
		log.Fatalf("\n\nERROR: %v", err)
	}
}

func getLatestVersion(ctx context.Context) (string, error) {
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, "3lvia", "cli")
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(*release.TagName, "v"), nil
}
