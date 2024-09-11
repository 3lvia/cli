package build

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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
			Name:    "system-name",
			Aliases: []string{"s"},
			Usage:   "The system name to use",
			EnvVars: []string{"3LV_SYSTEM_NAME"},
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

	// Optional args
	systemName := c.String("system-name")
	buildContext := c.String("build-context")
	registry := c.String("registry")
	goMainPackageDirectory := c.String("go-main-package-directory")
	cacheTag := c.String("cache-tag")
	severity := c.String("severity")
	additionalTags := removeZeroValues(c.StringSlice("additional-tags"))
	includeFiles := removeZeroValues(c.StringSlice("include-files"))
	includeDirectories := removeZeroValues(c.StringSlice("include-directories"))
	push := c.Bool("push")
	generateOnly := c.Bool("generate-only")

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
		DockerfilePath: dockerfilePath,
		BuildContext:   buildContext,
		CacheTag:       cacheTag,
		Registry:       registry,
		Severity:       severity,
		AdditionalTags: additionalTags,
		Push:           push,
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

type BuildAndPushImageOptions struct {
	DockerfilePath string   // required
	BuildContext   string   // required
	CacheTag       string   // required
	Registry       string   // required
	Severity       string   // required
	AdditionalTags []string // required
	Push           bool     // required
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

	return append(append([]string{
		"buildx",
		"build",
		"-f",
		dockerfilePath,
		"--cache-to",
		"type=inline",
		"--cache-from",
		imageName + ":" + cacheTag,
	}, tagArguments...), buildContext)
}

func getImageName(
	systemName string,
	applicationName string,
	registry string,
) string {
	if systemName != "" {
		return registry + "/" + systemName + "-" + applicationName
	}
	return registry + "/" + applicationName
}

func buildAndPushImage(
	systemName string,
	applicationName string,
	options BuildAndPushImageOptions,
) error {
	imageName := getImageName(systemName, applicationName, options.Registry)

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

	err := scanImage(ScanImageOptions{
		ImageName: imageName + ":" + options.CacheTag,
		Severity:  options.Severity,
	})
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

func removeZeroValues(slice []string) []string {
	var result []string
	for _, value := range slice {
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}
