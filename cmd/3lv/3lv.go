package main

import (
	"log"
	"os"

	"github.com/3lvia/cli/pkg/build"
	"github.com/3lvia/cli/pkg/deploy"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "3lv",
		Usage:                "Command Line Interface tool for developing, building and deploying Elvia applications",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			build.Command,
			deploy.Command,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
