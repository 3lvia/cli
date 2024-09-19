package deploy

import (
	"fmt"
	"os"
	"os/exec"
)

type AuthenticateGKEOptions struct {
	GKEProjectID       string
	GKEClusterName     string
	GKEClusterLocation string
}

func authenticateGKE(environment string, options AuthenticateGKEOptions) error {
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
