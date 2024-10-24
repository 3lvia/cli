package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestGcloudGetCredentialsCommand1(t *testing.T) {
	const gkeClusterName = "my-sick-cluster"
	const gkeClusterLocation = "europe-west1"
	const gkeProjectID = "my-cool-project"
	const environment = "this-will-not-be-used"

	expectedCommandString := strings.Join(
		[]string{
			"gcloud",
			"container",
			"clusters",
			"get-credentials",
			gkeClusterName,
			"--region",
			gkeClusterLocation,
			"--project",
			gkeProjectID,
		},
		" ",
	)

	actualCommand := gcloudGetCredentialsCommand(
		environment,
		GcloudGetCredentialsCommandOptions{
			GKEClusterName:     gkeClusterName,
			GKEClusterLocation: gkeClusterLocation,
			GKEProjectID:       gkeProjectID,
			RunOptions:         &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
