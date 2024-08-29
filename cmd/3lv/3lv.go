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
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Build the project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "project-file",
						Aliases: []string{"f"},
						Usage:   "The project file to use",
					},
					&cli.StringFlag{
						Name:    "build-context",
						Aliases: []string{"c"},
						Usage:   "The build context to use",
					},
					&cli.StringSliceFlag{
						Name:    "include-files",
						Aliases: []string{"i"},
						Usage:   "The files to include in the build context",
					},
					&cli.StringSliceFlag{
						Name:    "include-directories",
						Aliases: []string{"I"},
						Usage:   "The directories to include in the build context",
					},
					&cli.StringFlag{
						Name:  "go-main-package-directory",
						Usage: "The main package directory to use when building a Go application",
					},
					&cli.StringFlag{
						Name:    "system-name",
						Aliases: []string{"s"},
						Usage:   "The system name to use",
					},
					&cli.StringFlag{
						Name:    "registry",
						Aliases: []string{"r"},
						Usage:   "The registry to use",
						Value:   "acr",
					},
					&cli.BoolFlag{
						Name:    "push",
						Aliases: []string{"p"},
						Usage:   "Push the image to the registry",
						Value:   false,
					},
				},
				Action: build.Build,
			},
			{
				Name:    "deploy",
				Aliases: []string{"d"},
				Usage:   "Deploy the project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "system-name",
						Aliases: []string{"s"},
						Usage:   "The system name to use",
					},
					&cli.StringFlag{
						Name:    "environment",
						Aliases: []string{"e"},
						Usage:   "The environment to deploy to",
						Value:   "dev",
					},
					&cli.StringFlag{
						Name:    "helm-values-file",
						Aliases: []string{"v"},
						Usage:   "The helm values file to use",
					},
					&cli.StringFlag{
						Name:    "workload-type",
						Aliases: []string{"w"},
						Usage:   "The workload type to use",
						Value:   "deployment",
					},
					&cli.StringFlag{
						Name:    "runtime-cloud-provider",
						Aliases: []string{"r"},
						Usage:   "The runtime cloud provider to use",
						Value:   "aks",
					},
				},
				Action: deploy.Deploy,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
