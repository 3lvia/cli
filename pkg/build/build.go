package build

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/scan"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

const commandName = "build"

var Command *cli.Command = &cli.Command{
	Name:    commandName,
	Aliases: []string{"b"},
	Usage:   "Build the project",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "project-file",
			Aliases:  []string{"f"},
			Usage:    "The project file to use",
			Required: true,
			EnvVars:  []string{"3LV_PROJECT_FILE"},
		},
		&cli.StringFlag{
			Name:     "system-name",
			Aliases:  []string{"s"},
			Usage:    "The system name to use",
			Required: true,
			EnvVars:  []string{"3LV_SYSTEM_NAME"},
		},
		&cli.StringFlag{
			Name:    "build-context",
			Aliases: []string{"c"},
			Usage:   "The build context to use",
			EnvVars: []string{"3LV_BUILD_CONTEXT"},
		},
		&cli.StringFlag{
			Name:    "registry",
			Aliases: []string{"r"},
			Usage:   "The registry to use. Image name will be prefixed with this value.",
			Value:   "containerregistryelvia.azurecr.io",
			EnvVars: []string{"3LV_REGISTRY"},
		},
		&cli.StringFlag{
			Name:    "go-main-package-directory",
			Usage:   "The main package directory to use when building a Go application",
			EnvVars: []string{"3LV_GO_MAIN_PACKAGE_DIRECTORY"},
		},
		&cli.StringFlag{
			Name:    "cache-tag",
			Usage:   "The cache tag to use",
			Value:   "latest-cache",
			EnvVars: []string{"3LV_CACHE_TAG"},
		},
		&cli.StringFlag{
			Name:    "severity",
			Aliases: []string{"S"},
			Usage:   "The severity to use when scanning the image: can be any combination of CRITICAL, HIGH, MEDIUM, LOW, or UNKNOWN separated by commas",
			Value:   "CRITICAL,HIGH",
			EnvVars: []string{"3LV_SEVERITY"},
		},
		&cli.StringSliceFlag{
			Name:    "additional-tags",
			Aliases: []string{"t"},
			Usage:   "The additional tags to use",
			EnvVars: []string{"3LV_ADDITIONAL_TAGS"},
		},
		&cli.StringSliceFlag{
			Name:    "include-files",
			Aliases: []string{"i"},
			Usage:   "The files to include in the build context",
			EnvVars: []string{"3LV_INCLUDE_FILES"},
		},
		&cli.StringSliceFlag{
			Name:    "include-directories",
			Aliases: []string{"I"},
			Usage:   "The directories to include in the build context",
			EnvVars: []string{"3LV_INCLUDE_DIRECTORIES"},
		},
		&cli.StringSliceFlag{
			Name:    "scan-formats",
			Aliases: []string{"F"},
			Usage:   "The formats to use when outputting the scan results: can be table, json, sarif or markdown.",
			Value:   cli.NewStringSlice("table"),
			Action: func(c *cli.Context, formats []string) error {
				for _, format := range formats {
					if format != "table" && format != "json" && format != "sarif" && format != "markdown" {
						return cli.Exit("Invalid format provided", 1)
					}
				}

				return nil
			},
			EnvVars: []string{"3LV_SCAN_FORMATS"},
		},
		&cli.BoolFlag{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Push the image to the registry",
			Value:   false,
			EnvVars: []string{"3LV_PUSH"},
		},
		&cli.BoolFlag{
			Name:    "generate-only",
			Aliases: []string{"G"},
			Usage:   "Generates a Dockerfile, but does not build the image",
			Value:   false,
			EnvVars: []string{"3LV_GENERATE_ONLY"},
		},
		&cli.BoolFlag{
			Name:    "scan-disable-error",
			Aliases: []string{"D"},
			Usage:   "Disables Trivy scan returning a non-zero exit code if vulnerabilities are found",
			Value:   false,
			EnvVars: []string{"3LV_SCAN_DISABLE_ERROR"},
		},
		&cli.BoolFlag{
			Name:    "scan-skip-db-update",
			Aliases: []string{"U"},
			Usage:   "Skip updating the Trivy database.",
			Value:   false,
			EnvVars: []string{"3LV_SCAN_SKIP_DB_UPDATE"},
		},
	},
	Action: Build,
}

func Build(c *cli.Context) error {
	if c.NArg() <= 0 {
		return cli.ShowCommandHelp(c, commandName)
	}

	// Required args
	applicationName := c.Args().First()
	if applicationName == "" {
		log.Println("Application name not provided")
		return cli.ShowCommandHelp(c, commandName)
	}
	projectFile := c.String("project-file")
	if projectFile == "" {
		log.Println("Project file not provided")
		return cli.ShowCommandHelp(c, commandName)
	}
	systemName := c.String("system-name")
	if systemName == "" {
		log.Println("System name not provided")
		return cli.ShowCommandHelp(c, commandName)
	}

	generateOptions := GenerateDockerfileOptions{
		GoMainPackageDirectory: c.String("go-main-package-directory"),
		BuildContext:           c.String("build-context"),
		IncludeFiles:           utils.RemoveZeroValues(c.StringSlice("include-files")),
		IncludeDirectories:     utils.RemoveZeroValues(c.StringSlice("include-directories")),
	}

	dockerfilePath, buildContext, err := generateDockerfile(
		projectFile,
		applicationName,
		generateOptions,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if c.Bool("generate-only") {
		log.Printf("Dockerfile generated at %s\n", dockerfilePath)
		return nil
	}

	cacheTag := c.String("cache-tag")
	registry := c.String("registry")

	imageName, err := getImageName(
		registry,
		systemName,
		applicationName,
	)

	buildImageCommandOutput := buildImageCommand(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		utils.RemoveZeroValues(c.StringSlice("additional-tags")),
		nil,
	)
	if command.IsError(buildImageCommandOutput) {
		return cli.Exit(buildImageCommandOutput.Error, 1)
	}

	scanErr := scan.ScanImage(
		imageName+":"+cacheTag,
		c.String("severity"),
		utils.RemoveZeroValues(c.StringSlice("scan-formats")),
		c.Bool("scan-disable-error"),
		c.Bool("scan-skip-db-update"),
	)

	push := c.Bool("push")

	if push && scanErr != nil {
		pushImageOutput := pushImageCommand(
			imageName,
			cacheTag,
			false,
			nil,
		)

		if command.IsError(pushImageOutput) {
			return fmt.Errorf(
				"Failed to push Docker image cache to tag %s after scan reported vulnerabilities: %w",
				cacheTag,
				err,
			)
		}
	}

	if scanErr != nil {
		return scanErr
	}

	if push {
		pushImageOutput := pushImageCommand(
			imageName,
			cacheTag,
			true,
			nil,
		)

		if command.IsError(pushImageOutput) {
			return fmt.Errorf("Failed to push Docker image: %w", err)
		}
	}

	return nil
}

func getImageName(
	registry string,
	systemName string,
	applicationName string,
) (string, error) {
	if registry == "" {
		return "", fmt.Errorf("getImageName: Registry not provided")
	}
	if systemName == "" {
		return "", fmt.Errorf("getImageName: System name not provided")
	}
	if applicationName == "" {
		return "", fmt.Errorf("getImageName: Application name not provided")
	}

	if strings.Contains(registry, "azurecr.io") || strings.Contains(registry, "gcr.io") {
		return strings.ToLower(fmt.Sprintf("%s/%s-%s", registry, systemName, applicationName)), nil
	}
	return strings.ToLower(fmt.Sprintf("%s/%s/%s", registry, systemName, applicationName)), nil
}

func buildImageCommand(
	dockerfilePath string,
	buildContext string,
	imageName string,
	cacheTag string,
	additionalTags []string,
	options *command.RunOptions,
) command.Output {
	tags := func() []string {
		if len(additionalTags) == 0 {
			return []string{cacheTag}
		}

		return append(additionalTags, cacheTag)
	}()

	var tagArguments []string
	for _, tag := range tags {
		tagArguments = append(tagArguments, "-t")
		tagArguments = append(tagArguments, imageName+":"+tag)
	}

	buildCmd := exec.Command(
		"docker",
		"buildx",
		"build",
		"-f",
		dockerfilePath,
		"--load",
		"--cache-to",
		"type=inline",
		"--cache-from",
		imageName+":"+cacheTag,
	)

	buildCmd.Args = append(buildCmd.Args, tagArguments...)
	buildCmd.Args = append(buildCmd.Args, buildContext)

	return command.Run(*buildCmd, options)
}

func pushImageCommand(
	imageName string,
	cacheTag string,
	allTags bool,
	options *command.RunOptions,
) command.Output {
	if allTags {
		return command.Run(
			*exec.Command(
				"docker",
				"push",
				imageName,
				"--all-tags",
			),
			options,
		)
	}

	return command.Run(
		*exec.Command(
			"docker",
			"push",
			imageName+":"+cacheTag,
		),
		options,
	)
}
