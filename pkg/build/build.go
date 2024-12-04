package build

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/scan"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/urfave/cli/v3"
)

const commandName = "build"

var Command *cli.Command = &cli.Command{
	Name:    commandName,
	Aliases: []string{"b"},
	Usage:   "Build a Docker image from a project file.",
	Flags: []cli.Flag{
		shared.ProjectFileFlag(),
		shared.SystemNameFlag(
			"The system name to prefix the image name with. If not provided, we will try to use the current git repository name.",
			false,
		),
		shared.SeverityFlag("scan-severity"),
		shared.FormatsFlag("scan-formats"),
		shared.DisableErrorFlag("scan-disable-error"),
		shared.RegistryFlag("The container registry to use. Image name will be prefixed with this value."),
		&cli.StringFlag{
			Name:    "build-context",
			Aliases: []string{"c"},
			Usage:   "The directory to use as the build context for Docker, i.e. what files Docker will know about when building. We default to the directory of the project file. This means that if you need files outside of the directory of the project file, you need to specify this flag.",
			Sources: cli.EnvVars("3LV_BUILD_CONTEXT"),
		},
		&cli.StringFlag{
			Name:    "go-main-package-directory",
			Usage:   "The main package directory to use when building a Go application.",
			Sources: cli.EnvVars("3LV_GO_MAIN_PACKAGE_DIRECTORY"),
		},
		&cli.StringFlag{
			Name:    "cache-tag",
			Usage:   "The tag to use for the cache image.",
			Value:   "latest-cache",
			Sources: cli.EnvVars("3LV_CACHE_TAG"),
		},
		&cli.StringFlag{
			Name:    "azure-tenant-id",
			Usage:   "The tenant ID to use when authenticating with the Azure Container Registry.",
			Hidden:  true,
			Sources: cli.EnvVars("3LV_AZURE_TENANT_ID"),
		},
		&cli.StringFlag{
			Name:    "azure-subscription-id",
			Usage:   "The subscription ID to use when authenticating with the Azure Container Registry.",
			Hidden:  true,
			Sources: cli.EnvVars("3LV_AZURE_SUBSCRIPTION_ID"),
		},
		&cli.StringFlag{
			Name:    "azure-client-id",
			Usage:   "The client ID to use when authenticating with the Azure Container registry. Must be combined with --azure-federated-token.",
			Hidden:  true,
			Sources: cli.EnvVars("3LV_AZURE_CLIENT_ID"),
		},
		&cli.StringFlag{
			Name:    "azure-federated-token",
			Usage:   "The federated token to use when authenticating with the Azure Container Registry. Must be combined with --client-id.",
			Hidden:  true,
			Sources: cli.EnvVars("3LV_AZURE_FEDERATED_TOKEN"),
		},
		&cli.StringSliceFlag{
			Name:    "additional-tags",
			Aliases: []string{"t"},
			Usage:   "Additional tags to use when pushing the image to the registry.",
			Sources: cli.EnvVars("3LV_ADDITIONAL_TAGS"),
		},
		&cli.StringSliceFlag{
			Name:    "include-files",
			Aliases: []string{"i"},
			Usage:   "A list of files to include in the Docker image. Currently only supported for Go applications.",
			Sources: cli.EnvVars("3LV_INCLUDE_FILES"),
		},
		&cli.StringSliceFlag{
			Name:    "include-directories",
			Aliases: []string{"I"},
			Usage:   "A list of directories to include in the Docker image. Currently only supported for Go applications.",
			Sources: cli.EnvVars("3LV_INCLUDE_DIRECTORIES"),
		},
		&cli.BoolFlag{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Push the image to the registry.",
			Value:   false,
			Sources: cli.EnvVars("3LV_PUSH"),
		},
		&cli.BoolFlag{
			Name:    "generate-only",
			Aliases: []string{"G"},
			Usage:   "Generates a Dockerfile, but does not build the image.",
			Value:   false,
			Sources: cli.EnvVars("3LV_GENERATE_ONLY"),
		},
		&cli.BoolFlag{
			Name:    "skip-authentication",
			Usage:   "Skip authentication before pushing the image to the registry.",
			Value:   false,
			Sources: cli.EnvVars("3LV_SKIP_AUTHENTICATION"),
		},
	},
	Action: Build,
}

func Build(ctx context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		return cli.ShowAppHelp(c)
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
	systemName, err := func() (string, error) {
		possibleSystemName := c.String("system-name")

		if possibleSystemName == "" {
			style.Print(
				"System name not provided, will try to use the current git repository name.",
				nil,
			)

			repositoryName, err := utils.ResolveRepositoryName("")
			if err != nil {
				return "", err
			}

			return repositoryName, nil
		}

		return possibleSystemName, nil
	}()
	if err != nil {
		return cli.Exit(err, 1)
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
		style.Print(
			fmt.Sprintf("Dockerfile generated at %s\n", dockerfilePath),
			nil,
		)
		return nil
	}

	cacheTag := c.String("cache-tag")
	registry := utils.StringWithDefault(c.String("registry"), "containerregistryelvia.azurecr.io")

	push := c.Bool("push")
	skipAuthentication := c.Bool("skip-authentication") || !push

	if strings.Contains(registry, "azurecr.io") && !skipAuthentication {
		style.Print("Azure registry detected, will try to authenticate with Azure.", nil)

		azureTenantID := utils.StringWithDefault(
			c.String("azure-tenant-id"),
			auth.ElviaTenantID,
		)
		azureSubscriptionID := utils.StringWithDefault(
			c.String("azure-subscription-id"),
			auth.ElviaDefaultRuntimeSubscriptionID,
		)

		options := &auth.AzLoginCommandOptions{
			ClientID:       c.String("azure-client-id"),
			FederatedToken: c.String("azure-federated-token"),
		}

		err := auth.AuthenticateAzure(
			azureTenantID,
			azureSubscriptionID,
			options,
		)
		if err != nil {
			return cli.Exit(err, 1)
		}

		registryName, err := func() (string, error) {
			split := strings.Split(registry, ".")
			if len(split) <= 0 {
				return "", fmt.Errorf("Invalid registry name: %s", registry)
			}
			return split[0], nil
		}()
		if err != nil {
			return cli.Exit(err, 1)
		}

		azAcrLoginCommandOutput := azAcrLoginCommand(
			registryName,
			nil,
		)
		if command.IsError(azAcrLoginCommandOutput) {
			return cli.Exit(
				fmt.Errorf(
					"Failed to authenticate to Azure Container Registry: %w",
					azAcrLoginCommandOutput.Error,
				),
				1,
			)
		}

	}

	imageName, err := GetImageName(
		registry,
		systemName,
		applicationName,
	)

	additionalTags := utils.RemoveZeroValues(c.StringSlice("additional-tags"))

	buildImageCommandOutput := buildImageCommand(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		additionalTags,
		nil,
	)
	if command.IsError(buildImageCommandOutput) {
		return cli.Exit(buildImageCommandOutput.Error, 1)
	}

	scanErr := scan.ScanImage(
		imageName+":"+cacheTag,
		c.String("scan-severity"),
		utils.RemoveZeroValues(c.StringSlice("scan-formats")),
		c.Bool("scan-disable-error"),
	)

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
			return fmt.Errorf("Failed to push Docker image. If using GHCR, please login using the command `gh auth login` first. %w", err)
		}
	}

	outputDirectory := os.TempDir() + "/3lv-cli-output"
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := os.Mkdir(outputDirectory, 0700)
		if err != nil {
			return cli.Exit(err, 1)
		}
	}

	firstAdditionalTagThatsNotCacheTag := func() string {
		for _, tag := range additionalTags {
			if tag != cacheTag {
				return tag
			}
		}

		return cacheTag
	}()

	err = os.WriteFile(
		outputDirectory+"/image-name",
		[]byte(imageName+":"+firstAdditionalTagThatsNotCacheTag),
		0700,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	return nil
}

func GetImageName(
	registry string,
	systemName string,
	applicationName string,
) (string, error) {
	if registry == "" {
		return "", fmt.Errorf("Registry not provided")
	}
	if systemName == "" {
		return "", fmt.Errorf("System name not provided")
	}
	if applicationName == "" {
		return "", fmt.Errorf("Application name not provided")
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

func azAcrLoginCommand(
	registryName string,
	options *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"az",
			"acr",
			"login",
			"--name",
			registryName,
		),
		options,
	)
}
