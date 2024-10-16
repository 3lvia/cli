package deploy

import (
	"fmt"
	"os/exec"

	"github.com/3lvia/cli/pkg/command"
)

type SetupGKEOptions struct {
	GKEProjectID       string
	GKEClusterName     string
	GKEClusterLocation string
}

func setupGKE(
	environment string,
	skipAuthentication bool,
	options SetupGKEOptions,
) error {
	if !skipAuthentication {
		gcloudAuthLoginOutput := gcloudAuthLoginCommand(nil)
		if command.IsError(gcloudAuthLoginOutput) {
			return fmt.Errorf("Failed to authenticate to GKE: %w", gcloudAuthLoginOutput.Error)
		}
	}

	gcloudGetCredentialsOutput := gcloudGetCredentialsCommand(
		environment,
		GcloudGetCredentialsCommandOptions{
			GKEProjectID:       options.GKEProjectID,
			GKEClusterName:     options.GKEClusterName,
			GKEClusterLocation: options.GKEClusterLocation,
			RunOptions:         nil,
		},
	)

	if command.IsError(gcloudGetCredentialsOutput) {
		return fmt.Errorf("Failed to get GKE credentials: %w", gcloudGetCredentialsOutput.Error)
	}

	return nil
}

func gcloudAuthLoginCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command(
			"gcloud",
			"auth",
			"login",
		),
		runOptions,
	)
}

type GcloudGetCredentialsCommandOptions struct {
	GKEProjectID       string
	GKEClusterName     string
	GKEClusterLocation string
	RunOptions         *command.RunOptions
}

func gcloudGetCredentialsCommand(
	environment string,
	options GcloudGetCredentialsCommandOptions,
) command.Output {
	if environment == "" &&
		(options.GKEProjectID == "" ||
			options.GKEClusterName == "" ||
			options.GKEClusterLocation == "") {
		return command.ErrorString("environment must be set if any of the GKE options are not set")
	}

	gkeProjectID := func() string {
		if options.GKEProjectID == "" {
			return "elvia-runtimeservice-" + environment
		}

		return options.GKEProjectID
	}()

	gkeClusterName := func() string {
		if options.GKEClusterName == "" {
			return "runtimeservice-gke-" + environment
		}

		return options.GKEClusterName
	}()

	gkeClusterLocation := func() string {
		if options.GKEClusterLocation == "" {
			return "europe-west1"
		}

		return options.GKEClusterLocation
	}()

	return command.Run(
		*exec.Command(
			"gcloud",
			"container",
			"clusters",
			"get-credentials",
			gkeClusterName,
			"--region",
			gkeClusterLocation,
			"--project",
			gkeProjectID,
		),
		options.RunOptions,
	)
}
