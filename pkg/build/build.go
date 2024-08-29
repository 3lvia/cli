package build

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

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
	includeFiles := c.StringSlice("include-files")
	includeDirectories := c.StringSlice("include-directories")
	systemName := c.String("system-name")
	registry := c.String("registry")
	push := c.Bool("push")

	generateOptions := GenerateDockerfileOptions{
		ProjectFile:        projectFile,
		BuildContext:       buildContext,
		IncludeFiles:       includeFiles,
		IncludeDirectories: includeDirectories,
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
	ApplicationName string
	SystemName      string
	Registry        string
	Tags            []string
	Push            bool
}

func buildAndPushImage(options BuildAndPushImageOptions) error {
	registry := getRegistry(options.Registry)
	imageName := registry + "/" + options.SystemName + "-" + options.ApplicationName

	var tagArguments []string
	if len(options.Tags) == 0 {
		options.Tags = []string{"latest-cache", "latest"}
	}
	for _, tag := range options.Tags {
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

	err := scanImage(imageName)
	if err != nil {
		return err
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

func scanImage(imageName string) error {
	scanCmd := exec.Command(
		"trivy",
		"image",
		"--severity",
		"CRITICAL,HIGH",
		imageName,
	)
	scanCmd.Stdout = os.Stdout
	scanCmd.Stderr = os.Stderr

	if err := scanCmd.Run(); err != nil {
		return fmt.Errorf("Failed to scan Docker image: %w", err)
	}

	return nil
}
