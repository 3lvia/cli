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

func checkVersionAfter(ctx context.Context, c *cli.Command) error {
	latestVersion, _ := upgrade.GetLatestCLIVersion(ctx)
	if semver.Compare("v"+c.Version, "v"+latestVersion) == -1 {
		style.Print(
			fmt.Sprintf("\n\nA new version of 3lv is available! %s -> %s", c.Version, latestVersion),
			&style.PrintOptions{Color: "yellow"},
		)

		style.Print(
			"Run `3lv upgrade` to update to the latest version.",
			&style.PrintOptions{Color: "yellow"},
		)
	}

	return nil
}

func withCheckVersionAfter(c *cli.Command) *cli.Command {
	c.After = func(ctx context.Context, c_ *cli.Command) error {
		return checkVersionAfter(ctx, c_)
	}

	return c
}

func WithCheckVersionAfterCommands(commands []*cli.Command) []*cli.Command {
	return lo.Map(commands, func(c *cli.Command, _ int) *cli.Command {
		return withCheckVersionAfter(c)
	})
}
