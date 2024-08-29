package main

import (
	"embed"
	"log"
	"os"

	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/deploy"
	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version embed.FS

func main() {
	versionFile, err := version.ReadFile("VERSION")
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:                 "3lv",
		Usage:                "Command Line Interface tool for developing, building and deploying Elvia applications",
		EnableBashCompletion: true,
		Version:              string(versionFile),
		Commands: []*cli.Command{
			build.Command,
			deploy.Command,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
