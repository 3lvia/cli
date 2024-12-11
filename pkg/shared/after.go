package shared

import (
	"context"
	"fmt"

	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/upgrade"
	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
	"golang.org/x/mod/semver"
)

func checkVersionAfter(ctx context.Context, version string) error {
	latestVersion, _ := upgrade.GetLatestCLIVersion(ctx)
	if semver.Compare("v"+version, "v"+latestVersion) == -1 {
		style.Print(
			fmt.Sprintf("\n\nA new version of 3lv is available! %s -> %s", version, latestVersion),
			&style.PrintOptions{Color: "yellow"},
		)

		style.Print(
			"Run `3lv upgrade` to update to the latest version.",
			&style.PrintOptions{Color: "yellow"},
		)
	}

	return nil
}

func withCheckVersionAfter(c *cli.Command, version string) *cli.Command {
	c.After = func(ctx context.Context, _ *cli.Command) error {
		return checkVersionAfter(ctx, version)
	}

	return c
}

func WithCheckVersionAfterCommands(commands []*cli.Command, version string) []*cli.Command {
	return lo.Map(commands, func(c *cli.Command, _ int) *cli.Command {
		return withCheckVersionAfter(c, version)
	})
}
