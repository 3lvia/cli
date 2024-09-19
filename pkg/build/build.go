package build

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/scan"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v2"
)

var Command *cli.Command = &cli.Command{
	Name:    "build",
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
		return cli.Exit("No input provided", 1)
	}

	// Required args
	applicationName := c.Args().First()
	if applicationName == "" {
		return cli.Exit("Application name not provided", 1)
	}

	projectFile := c.String("project-file")
	if projectFile == "" {
		return cli.Exit("Project file not provided", 1)
	}

	systemName := c.String("system-name")
	if systemName == "" {
		return cli.Exit("System name not provided", 1)
	}

	// Optional args
	buildContext := c.String("build-context")
	registry := c.String("registry")
	goMainPackageDirectory := c.String("go-main-package-directory")
	cacheTag := c.String("cache-tag")
	severity := c.String("severity")
	additionalTags := utils.RemoveZeroValues(c.StringSlice("additional-tags"))
	includeFiles := utils.RemoveZeroValues(c.StringSlice("include-files"))
	includeDirectories := utils.RemoveZeroValues(c.StringSlice("include-directories"))
	scanFormats := utils.RemoveZeroValues(c.StringSlice("scan-formats"))
	scanSkipDBUpdate := c.Bool("scan-skip-db-update")
	push := c.Bool("push")
	generateOnly := c.Bool("generate-only")
	scanDisableError := c.Bool("scan-disable-error")

	generateOptions := GenerateDockerfileOptions{
		GoMainPackageDirectory: goMainPackageDirectory,
		BuildContext:           buildContext,
		IncludeFiles:           includeFiles,
		IncludeDirectories:     includeDirectories,
	}

	dockerfilePath, buildContext, err := generateDockerfile(
		projectFile,
		applicationName,
		generateOptions,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if generateOnly {
		log.Printf("Dockerfile generated at %s\n", dockerfilePath)
		return nil
	}

	buildOptions := BuildAndPushImageOptions{
		DockerfilePath:   dockerfilePath,
		BuildContext:     buildContext,
		CacheTag:         cacheTag,
		Registry:         registry,
		Severity:         severity,
		ScanFormats:      scanFormats,
		AdditionalTags:   additionalTags,
		Push:             push,
		ScanSkipDBUpdate: scanSkipDBUpdate,
		ScanDisableError: scanDisableError,
	}

	err = buildAndPushImage(
		systemName,
		applicationName,
		buildOptions,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func constructBuildCommandArguments(
	dockerfilePath string,
	buildContext string,
	imageName string,
	cacheTag string,
	additionalTags []string,
) []string {
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

	return append(
		append(
			[]string{
				"buildx",
				"build",
				"-f",
				dockerfilePath,
				"--load",
				"--cache-to",
				"type=inline",
				"--cache-from",
				imageName + ":" + cacheTag,
			},
			tagArguments...,
		),
		buildContext,
	)
}

func getImageName(
	registry string,
	systemName string,
	applicationName string,
) string {
	if strings.Contains(registry, "azurecr.io") {
		return strings.ToLower(fmt.Sprintf("%s/%s-%s", registry, systemName, applicationName))
	}
	return strings.ToLower(fmt.Sprintf("%s/%s/%s", registry, systemName, applicationName))
}

type BuildAndPushImageOptions struct {
	DockerfilePath   string   // required
	BuildContext     string   // required
	CacheTag         string   // required
	Registry         string   // required
	Severity         string   // required
	ScanFormats      []string // required
	AdditionalTags   []string // required
	Push             bool     // required
	ScanDisableError bool     // required
	ScanSkipDBUpdate bool     // required
}

func buildAndPushImage(
	systemName string,
	applicationName string,
	options BuildAndPushImageOptions,
) error {
	imageName := getImageName(options.Registry, systemName, applicationName)

	buildCmd := exec.Command(
		"docker",
		constructBuildCommandArguments(
			options.DockerfilePath,
			options.BuildContext,
			imageName,
			options.CacheTag,
			options.AdditionalTags,
		)...,
	)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	log.Println(buildCmd.String())

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("Failed to build Docker image: %w", err)
	}

	err := scan.ScanImage(
		imageName+":"+options.CacheTag,
		options.Severity,
		options.ScanFormats,
		options.ScanDisableError,
		options.ScanSkipDBUpdate,
	)
	if err != nil {
		return err
	}

	if options.Push {
		pushCmd := exec.Command(
			"docker",
			"push",
			imageName,
			"--all-tags",
		)
		pushCmd.Stdout = os.Stdout
		pushCmd.Stderr = os.Stderr

		if err := pushCmd.Run(); err != nil {
			return fmt.Errorf("Failed to push Docker image: %w", err)
		}
	}

	return nil
}
