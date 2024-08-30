package build

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"

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
			Usage:   "The registry to use",
			Value:   "acr",
			Action: func(c *cli.Context, registry string) error {
				allowedRegistries := []string{"acr", "ghcr"}
				if !slices.Contains(allowedRegistries, registry) {
					return fmt.Errorf("Unknown registry '%s'; allowed values are %v", registry, allowedRegistries)
				}

				return nil
			},
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
	additionalTags := c.StringSlice("additional-tags")
	includeFiles := c.StringSlice("include-files")
	includeDirectories := c.StringSlice("include-directories")
	push := c.Bool("push")

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
	tags []string,
) []string {
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
	}, tagArguments...), buildContext)
}

func buildAndPushImage(
	systemName string,
	applicationName string,
	options BuildAndPushImageOptions,
) error {
	registry := getRegistry(options.Registry)
	imageName := registry + "/" + systemName + "-" + applicationName

	tags := func() []string {
		if len(options.AdditionalTags) == 0 {
			return []string{options.CacheTag}
		}

		return append(options.AdditionalTags, options.CacheTag)
	}()

	buildCmd := exec.Command(
		"docker",
		constructBuildCommandArguments(
			options.DockerfilePath,
			options.BuildContext,
			imageName,
			tags,
		)...,
	)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	log.Println(buildCmd.String())

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("Failed to build Docker image: %w", err)
	}

	err := scanImage(ScanImageOptions{
		ImageName: imageName + ":" + tags[0],
		Severity:  options.Severity,
	})
	if err != nil {
		return err
	}

	if options.Push {
		/*
			err := loginToRegistry(options.Registry)
			if err != nil {
				return err
			}
		*/

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

/*
func loginToRegistry(registry string) error {
	switch registry {
	case "acr":
		loginCmd := exec.Command(
			"az",
			"acr",
			"login",
			"--name",
			"containerregistryelvia",
		)
		loginCmd.Stdout = os.Stdout
		loginCmd.Stderr = os.Stderr

		if err := loginCmd.Run(); err != nil {
			return fmt.Errorf("Failed to login to Azure Container Registry: %w", err)
		}

		return nil
	case "ghcr":
		return fmt.Errorf("Not implemented")
	default:
		return fmt.Errorf("Unknown registry")
	}
}
*/

func getRegistry(registry string) string {
	if registry == "" || registry == "acr" {
		return "containerregistryelvia.azurecr.io"
	} else if registry == "ghcr" {
		return "ghcr.io/3lvia"
	}

	return registry
}

type ScanImageOptions struct {
	ImageName string // required
	Severity  string // required
}

func scanImage(options ScanImageOptions) error {
	scanCmd := exec.Command(
		"trivy",
		"image",
		"--severity",
		options.Severity,
		options.ImageName,
	)
	scanCmd.Stdout = os.Stdout
	scanCmd.Stderr = os.Stderr

	if err := scanCmd.Run(); err != nil {
		return fmt.Errorf("Failed to scan Docker image: %w", err)
	}

	return nil
}
