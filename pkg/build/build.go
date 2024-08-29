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
			Name:    "project-file",
			Aliases: []string{"f"},
			Usage:   "The project file to use",
		},
		&cli.StringFlag{
			Name:    "build-context",
			Aliases: []string{"c"},
			Usage:   "The build context to use",
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
		&cli.StringFlag{
			Name:  "go-main-package-directory",
			Usage: "The main package directory to use when building a Go application",
		},
		&cli.StringFlag{
			Name:  "cache-tag",
			Usage: "The cache tag to use",
			Value: "latest-cache",
		},
		&cli.StringFlag{
			Name:    "severity",
			Aliases: []string{"S"},
			Usage:   "The severity to use when scanning the image: can be any combination of CRITICAL, HIGH, MEDIUM, LOW, or UNKNOWN separated by commas",
			Value:   "CRITICAL,HIGH",
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
		&cli.BoolFlag{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Push the image to the registry",
			Value:   false,
		},
	},
	Action: Build,
}

func Build(c *cli.Context) error {
	if c.NArg() <= 0 {
		return cli.Exit("No input provided", 1)
	}

	applicationName := c.Args().First()
	if applicationName == "" {
		return cli.Exit("Application name not provided", 1)
	}

	projectFile := c.String("project-file")
	if projectFile == "" {
		return cli.Exit("Project file not provided", 1)
	}
	buildContext := c.String("build-context")
	goMainPackageDirectory := c.String("go-main-package-directory")
	systemName := c.String("system-name")
	registry := c.String("registry")
	severity := c.String("severity")
	includeFiles := c.StringSlice("include-files")
	includeDirectories := c.StringSlice("include-directories")
	push := c.Bool("push")

	generateOptions := GenerateDockerfileOptions{
		ProjectFile:            projectFile,
		ApplicationName:        applicationName,
		GoMainPackageDirectory: goMainPackageDirectory,
		BuildContext:           buildContext,
		IncludeFiles:           includeFiles,
		IncludeDirectories:     includeDirectories,
	}

	dockerfilePath, buildContext, err := generateDockerfile(generateOptions)
	if err != nil {
		return cli.Exit(err, 1)
	}

	buildOptions := BuildAndPushImageOptions{
		DockerfilePath:  dockerfilePath,
		BuildContext:    buildContext,
		ApplicationName: applicationName,
		SystemName:      systemName,
		Registry:        registry,
		Severity:        severity,
		Push:            push,
	}

	err = buildAndPushImage(buildOptions)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

type BuildAndPushImageOptions struct {
	DockerfilePath  string
	BuildContext    string
	CacheTag        string
	ApplicationName string
	SystemName      string
	Registry        string
	Severity        string
	AdditionalTags  []string
	Push            bool
}

func buildAndPushImage(options BuildAndPushImageOptions) error {
	registry := getRegistry(options.Registry)
	imageName := registry + "/" + options.SystemName + "-" + options.ApplicationName

	tags := func() []string {
		if len(options.AdditionalTags) == 0 {
			return []string{options.CacheTag}
		}

		return append(options.AdditionalTags, options.CacheTag)
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
		options.DockerfilePath,
	)
	buildCmd.Args = append(buildCmd.Args, tagArguments...)
	buildCmd.Args = append(buildCmd.Args, options.BuildContext)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	log.Println(buildCmd.String())

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("Failed to build Docker image: %w", err)
	}

	err := scanImage(ScanImageOptions{
		ImageName: imageName,
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
