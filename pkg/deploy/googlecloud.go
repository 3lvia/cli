package deploy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
		if err := authenticateGKE(); err != nil {
			return fmt.Errorf("Failed to authenticate to GKE: %w", err)
		}
	}

	if err := getGKECredentials(environment, GetGKECredentialsOptions(options)); err != nil {
		return fmt.Errorf("Failed to get GKE credentials: %w", err)
	}

	return nil
}

func authenticateGKE() error {
	gcloudLoginCmd := exec.Command(
		"gcloud",
		"auth",
		"login",
	)
	gcloudLoginCmd.Stdout = os.Stdout
	gcloudLoginCmd.Stderr = os.Stderr

	log.Print(gcloudLoginCmd.String())

	if err := gcloudLoginCmd.Run(); err != nil {
		return fmt.Errorf("Failed to authenticate to GKE: %w", err)
	}

	return nil
}

type GetGKECredentialsOptions struct {
	GKEProjectID       string
	GKEClusterName     string
	GKEClusterLocation string
}

func getGKECredentials(
	environment string,
	options GetGKECredentialsOptions,
) error {
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

	gcloudGetCredentialsCmd := exec.Command(
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
	gcloudGetCredentialsCmd.Stdout = os.Stdout
	gcloudGetCredentialsCmd.Stderr = os.Stderr

	if err := gcloudGetCredentialsCmd.Run(); err != nil {
		return fmt.Errorf("Failed to authenticate to GKE: %w", err)
	}

	return nil
}
