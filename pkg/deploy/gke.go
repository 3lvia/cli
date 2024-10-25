package deploy

import (
	"fmt"
	"os/exec"

	"github.com/3lvia/cli/pkg/auth"
	"github.com/3lvia/cli/pkg/command"
	"github.com/3lvia/cli/pkg/utils"
)

type SetupGKEOptions struct {
	ProjectID       string
	ClusterName     string
	ClusterLocation string
	UseInternalIP   bool
}

func setupGKE(
	environment string,
	skipAuthentication bool,
	options *SetupGKEOptions,
) error {
	if !skipAuthentication {
		err := auth.AuthenticateGoogle()
		if err != nil {
			return err
		}
	}

	gcloudGetCredentialsOutput := gcloudGetCredentialsCommand(
		environment,
		&GcloudGetCredentialsCommandOptions{
			ProjectID:       options.ProjectID,
			ClusterName:     options.ClusterName,
			ClusterLocation: options.ClusterLocation,
			RunOptions:      nil,
		},
	)

	if command.IsError(gcloudGetCredentialsOutput) {
		return fmt.Errorf("Failed to get GKE credentials: %w", gcloudGetCredentialsOutput.Error)
	}

	return nil
}

type GcloudGetCredentialsCommandOptions struct {
	ProjectID       string
	ClusterName     string
	ClusterLocation string
	UseInternalIP   bool
	RunOptions      *command.RunOptions
}

func gcloudGetCredentialsCommand(
	environment string,
	options *GcloudGetCredentialsCommandOptions,
) command.Output {
	if options == nil {
		options = &GcloudGetCredentialsCommandOptions{}
	}

	if environment == "" &&
		(options.ProjectID == "" ||
			options.ClusterName == "" ||
			options.ClusterLocation == "") {
		return command.ErrorString("environment must be set if any of the GKE options are not set")
	}

	gkeProjectID := func() string {
		if options.ProjectID == "" {
			return "elvia-runtimeservice-" + environment
		}

		return options.ProjectID
	}()

	gkeClusterName := func() string {
		if options.ClusterName == "" {
			return "runtimeservice-gke-" + environment
		}

		return options.ClusterName
	}()

	gkeClusterLocation := utils.StringWithDefault(options.ClusterLocation, "europe-west1")

	cmd := exec.Command(
		"gcloud",
		"container",
		"clusters",
		"get-credentials",
		gkeClusterName,
		"--region",
		gkeClusterLocation,
		"--project",
		gkeProjectID,
	)

	if options.UseInternalIP {
		cmd.Args = append(cmd.Args, "--internal-ip")
	}

	return command.Run(
		*cmd,
		options.RunOptions,
	)
}
