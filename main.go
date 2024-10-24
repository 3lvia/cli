package main

import (
	"embed"
	"log"
	"os"

	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/deploy"
	"github.com/3lvia/cli/pkg/scan"
	"github.com/urfave/cli/v2"
)

//go:embed VERSION
var version embed.FS

func main() {
	log.SetFlags(0)

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
			scan.Command,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("\n\nERROR: %v", err)
	}
}
