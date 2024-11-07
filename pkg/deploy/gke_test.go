package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestGcloudGetCredentialsCommand1(t *testing.T) {
	const clusterName = "my-sick-cluster"
	const clusterLocation = "europe-west1"
	const projectID = "my-cool-project"
	const environment = "this-will-not-be-used"

	expectedCommandString := strings.Join(
		[]string{
			"gcloud",
			"container",
			"clusters",
			"get-credentials",
			clusterName,
			"--region",
			clusterLocation,
			"--project",
			projectID,
		},
		" ",
	)

	actualCommand := gcloudGetCredentialsCommand(
		environment,
		&GcloudGetCredentialsCommandOptions{
			ClusterName:     clusterName,
			ClusterLocation: clusterLocation,
			ProjectID:       projectID,
			RunOptions:      &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestGcloudGetCredentialsCommand2(t *testing.T) {
	const clusterName = "my-sick-cluster"
	const clusterLocation = "europe-west1"
	const projectID = "my-cool-project"
	const environment = "this-will-not-be-used"

	expectedCommandString := strings.Join(
		[]string{
			"gcloud",
			"container",
			"clusters",
			"get-credentials",
			clusterName,
			"--region",
			clusterLocation,
			"--project",
			projectID,
		},
		" ",
	)

	actualCommand := gcloudGetCredentialsCommand(
		environment,
		&GcloudGetCredentialsCommandOptions{
			ClusterName:     clusterName,
			ClusterLocation: clusterLocation,
			ProjectID:       projectID,
			RunOptions:      &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestGcloudGetCredentialsCommand3(t *testing.T) {
	const clusterName = "my-sick-cluster"
	const clusterLocation = "europe-west1"
	const projectID = "my-cool-project"
	const environment = "this-will-not-be-used"

	expectedCommandString := strings.Join(
		[]string{
			"gcloud",
			"container",
			"clusters",
			"get-credentials",
			clusterName,
			"--region",
			clusterLocation,
			"--project",
			projectID,
		},
		" ",
	)

	actualCommand := gcloudGetCredentialsCommand(
		environment,
		&GcloudGetCredentialsCommandOptions{
			ClusterName:     clusterName,
			ClusterLocation: clusterLocation,
			ProjectID:       projectID,
			RunOptions:      &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestGcloudGetCredentialsCommand4(t *testing.T) {
	const environment = "dev"

	expectedCommandString := strings.Join(
		[]string{
			"gcloud",
			"container",
			"clusters",
			"get-credentials",
			"runtimeservice-gke-" + environment,
			"--region",
			"europe-west1",
			"--project",
			"elvia-runtimeservice-" + environment,
		},
		" ",
	)

	actualCommand := gcloudGetCredentialsCommand(
		environment,
		&GcloudGetCredentialsCommandOptions{
			RunOptions: &command.RunOptions{DryRun: true},
		},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
