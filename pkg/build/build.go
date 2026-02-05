package build

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/scan"
	"github.com/3lvia/cli/pkg/shared"
	"github.com/3lvia/cli/pkg/style"
	"github.com/3lvia/cli/pkg/utils"
	"github.com/samber/lo"
	"github.com/urfave/cli/v3"
)

const (
	commandName                   = "build"
	DefaultElviaContainerRegistry = "containerregistryelvia.azurecr.io"
	DefaultCacheTag               = "latest-cache"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:      commandName,
		Aliases:   []string{"b"},
		Usage:     "Build a Docker image from a project file.",
		UsageText: "3lv build [options] <application-name>",
		Flags: []cli.Flag{
			shared.ProjectFileFlag(),
			shared.SystemNameFlag(
				"The system name to prefix the image name with." +
					" If not provided, we will try to use the current git repository name.",
			),
			shared.SeverityFlag("scan-severity"),
			shared.FormatsFlag("scan-formats"),
			shared.DisableErrorFlag("scan-disable-error"),
			shared.RegistryFlag("The container registry to use. Image name will be prefixed with this value."),
			&cli.StringFlag{
				Name:    "build-context",
				Aliases: []string{"c"},
				Usage: "The directory to use as the Docker build context, i.e. what files Docker will know about when building." +
					" This means that if you need files outside of the directory of" +
					" the project file, you need to specify this flag.",
				Sources:     cli.EnvVars("3LV_BUILD_CONTEXT"),
				DefaultText: "directory of the project file",
			},
			&cli.StringFlag{
				Name:        "go-main-package-directory",
				Usage:       "The main package directory to use when building a Go application.",
				DefaultText: "\"./\"",
				Sources:     cli.EnvVars("3LV_GO_MAIN_PACKAGE_DIRECTORY"),
			},
			&cli.StringFlag{
				Name:    "cache-tag",
				Usage:   "The tag to use for the cache image.",
				Value:   DefaultCacheTag,
				Sources: cli.EnvVars("3LV_CACHE_TAG"),
			},
			&cli.BoolFlag{
				Name: "disable-cache",
				Usage: "Disable the use of cache when building the image." +
					" Cache produced by the build will still be pushed to the registry if --push is enabled.",
				Value:   false,
				Sources: cli.EnvVars("3LV_DISABLE_CACHE"),
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
				Name: "azure-client-id",
				Usage: "The client ID to use when authenticating with the Azure Container registry." +
					" Must be combined with --azure-federated-token.",
				Hidden:  true,
				Sources: cli.EnvVars("3LV_AZURE_CLIENT_ID"),
			},
			&cli.StringFlag{
				Name: "azure-federated-token",
				Usage: "The federated token to use when authenticating with the Azure Container Registry." +
					" Must be combined with --client-id.",
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
				Name: "build-args",
				Usage: "Build arguments to pass to the Docker build command." +
					" Should be in the format key=value, separated by commas.",
				Sources: cli.EnvVars("3LV_BUILD_ARGS"),
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
}

func Build(_ context.Context, c *cli.Command) error {
	if c.NArg() <= 0 {
		cli.ShowSubcommandHelpAndExit(c, 1)
	}

	applicationName := c.Args().First()
	if applicationName == "" {
		return cli.Exit("Application name not provided.", 1)
	}

	config, err := shared.GetConfig()
	if err != nil {
		style.PrintWarning(err.Error() + "\n")
	}

	configForApplication, err := config.GetConfigForApplication(applicationName)
	if !config.IsEmpty() && err != nil { // Ignore error if config is empty, will default to flags.
		style.PrintWarning(err.Error() + "\n")
	}

	projectFile := utils.FirstNonEmpty(c.String("project-file"), configForApplication.ProjectFile)
	if projectFile == "" {
		return cli.Exit("Project file not provided.", 1)
	}

	systemName := utils.FirstNonEmpty(c.String("system-name"), config.System)

	generateOptions := GenerateDockerfileOptions{
		GoMainPackageDirectory: utils.FirstNonEmpty(
			c.String("go-main-package-directory"),
			configForApplication.GoMainPackageDirectory,
			"./",
		),
		BuildContext: utils.FirstNonEmpty(c.String("build-context"), configForApplication.BuildContext),
	}

	dockerfilePath, buildContext, err := generateDockerfile(
		projectFile,
		generateOptions,
	)
	if err != nil {
		return cli.Exit(err, 1)
	}

	if c.Bool("generate-only") {
		newDockerfilePath, err := copyDockerfileToCurrentDirectory(
			dockerfilePath,
			c.Bool("non-interactive"),
		)
		if err != nil {
			return cli.Exit(err, 1)
		}

		style.PrintSuccess(fmt.Sprintf("Dockerfile generated at %s\n", newDockerfilePath))

		return nil
	}

	cacheTag := c.String("cache-tag")
	registry := utils.StringWithDefault(c.String("registry"), DefaultElviaContainerRegistry)

	push := c.Bool("push")
	skipAuthentication := c.Bool("skip-authentication") || !push

	if strings.Contains(registry, "azurecr.io") && !skipAuthentication {
		style.PrintInfo("Azure registry detected, will try to authenticate with Azure.")

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

		registryName, err := getRegistryName(registry)
		if err != nil {
			return cli.Exit(err, 1)
		}

		if azAcrLoginCommandOutput := azAcrLoginCommand(
			registryName,
			nil,
		); command.IsError(azAcrLoginCommandOutput) {
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
	if err != nil {
		return cli.Exit(err, 1)
	}

	disableCache := c.Bool("disable-cache")

	additionalTags := utils.RemoveZeroValues(c.StringSlice("additional-tags"))
	buildArgs := lo.SliceToMap(utils.RemoveZeroValues(c.StringSlice("build-args")), func(f string) (string, string) {
		split := strings.Split(f, "=")

		return split[0], split[1]
	})

	if buildImageCommandOutput := buildImageCommand(
		dockerfilePath,
		buildContext,
		imageName,
		cacheTag,
		disableCache,
		additionalTags,
		buildArgs,
		nil,
	); command.IsError(buildImageCommandOutput) {
		return cli.Exit(buildImageCommandOutput.Error, 1)
	}

	scanErr := scan.ScanImage(
		imageName+":"+cacheTag,
		c.String("scan-severity"),
		utils.RemoveZeroValues(c.StringSlice("scan-formats")),
		c.Bool("scan-disable-error"),
	)
	if push && scanErr != nil {
		if pushImageOutput := pushImageCommand(
			imageName,
			cacheTag,
			false,
			nil,
		); command.IsError(pushImageOutput) {
			return fmt.Errorf(
				"Failed to push Docker image cache to tag %s after scan reported vulnerabilities: %w",
				cacheTag,
				pushImageOutput.Error,
			)
		}
	}

	if scanErr != nil {
		return scanErr
	}

	if push {
		if pushImageOutput := pushImageCommand(
			imageName,
			cacheTag,
			true,
			nil,
		); command.IsError(pushImageOutput) {
			return fmt.Errorf(
				"Failed to push Docker image. If using GHCR, please login using the command `gh auth login` first. %w",
				pushImageOutput.Error,
			)
		}
	}

	outputDirectory := os.TempDir() + "/3lv-cli-output"
	if _, err := os.Stat(outputDirectory); os.IsNotExist(err) {
		err := os.Mkdir(outputDirectory, 0o700)
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
		0o700,
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
		return "", errors.New("Registry not provided")
	}

	if systemName == "" {
		return "", errors.New("System name not provided")
	}

	if applicationName == "" {
		return "", errors.New("Application name not provided")
	}

	return strings.ToLower(fmt.Sprintf("%s/%s/%s", registry, systemName, applicationName)), nil
}

func buildImageCommand(
	dockerfilePath string,
	buildContext string,
	imageName string,
	cacheTag string,
	disableCache bool,
	additionalTags []string,
	buildArgs map[string]string,
	options *command.RunOptions,
) command.Output {
	tags := func() []string {
		if len(additionalTags) == 0 {
			return []string{cacheTag}
		}

		return append(additionalTags, cacheTag)
	}()

	tagArguments := make([]string, 0, len(tags)*2)
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
	)

	if !disableCache {
		buildCmd.Args = append(buildCmd.Args,
			"--cache-from",
			imageName+":"+cacheTag,
		)
	}

	for key, value := range buildArgs {
		buildCmd.Args = append(buildCmd.Args, "--build-arg", key+"="+value)
	}

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

func getRegistryName(registry string) (string, error) {
	split := strings.Split(registry, ".")
	if len(split) <= 0 {
		return "", fmt.Errorf("Invalid registry name: %s", registry)
	}

	return split[0], nil
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

func copyDockerfileToCurrentDirectory(dockerfilePath string, nonInteractive bool) (string, error) {
	dockerfile, err := os.Open(dockerfilePath)
	if err != nil {
		return "", err
	}

	defer dockerfile.Close()

	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	newDockerfilePath := workingDirectory + "/Dockerfile"
	if _, err := os.Stat(newDockerfilePath); err == nil {
		overwriteFile, err := utils.PromptYesNo(
			"There is already a Dockerfile in the current directory. Do you want to overwrite it?",
			nonInteractive,
		)
		if err != nil {
			return "", err
		}

		if !overwriteFile {
			return "", errors.New("User chose not to overwrite existing Dockerfile")
		}
	}

	newDockerfile, err := os.Create(newDockerfilePath)
	if err != nil {
		return "", err
	}

	defer newDockerfile.Close()

	_, err = io.Copy(newDockerfile, dockerfile)
	if err != nil {
		return "", err
	}

	return newDockerfilePath, nil
}
