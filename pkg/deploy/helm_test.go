package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestCheckHelmInstalledCommand(t *testing.T) {
	expectedCommandString := "helm version"

	actualCommand := checkHelmInstalledCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmRepoAddCommand(t *testing.T) {
	expectedCommandString := strings.Join(
		[]string{
			"helm",
			"repo",
			"add",
			"elvia-charts",
			"https://raw.githubusercontent.com/3lvia/kubernetes-charts/master",
		},
		" ",
	)

	actualCommand := helmRepoAddCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmRepoUpdateCommand(t *testing.T) {
	expectedCommandString := "helm repo update"

	actualCommand := helmRepoUpdateCommand(
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmDeployCommand1(t *testing.T) {
	const systemName = "core"
	const helmValuesFile = ".github/deploy/values.yml"
	const applicationName = "demo-api"
	const environment = "dev"
	const workloadType = "deployment"
	const imageTag = "v12"
	const repositoryName = "core"
	const commitHash = "123456"

	expectedCommandString := strings.Join(
		[]string{
			"helm",
			"upgrade",
			"--debug",
			"--install",
			"-n",
			systemName,
			"-f",
			helmValuesFile,
			applicationName,
			"elvia-charts/elvia-" + workloadType,
			"--set-string",
			"environment=" + environment,
			"--set-string",
			"image.tag=" + imageTag,
			"--set-string",
			"labels.repositoryName=" + repositoryName,
			"--set-string",
			"labels.commitHash=\"" + commitHash + "\"",
		},
		" ",
	)

	actualCommand := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		repositoryName,
		commitHash,
		false,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmDeployCommand2(t *testing.T) {
	const systemName = "core"
	const helmValuesFile = ".github/deploy/values.yml"
	const applicationName = "demo-api"
	const environment = "prod"
	const workloadType = "statefulset"
	const imageTag = "v420"
	const repositoryName = "core-not-monorepo"
	const commitHash = "abcdef"

	expectedCommandString := strings.Join(
		[]string{
			"helm",
			"upgrade",
			"--debug",
			"--install",
			"-n",
			systemName,
			"-f",
			helmValuesFile,
			applicationName,
			"elvia-charts/elvia-" + workloadType,
			"--set-string",
			"environment=" + environment,
			"--set-string",
			"image.tag=" + imageTag,
			"--set-string",
			"labels.repositoryName=" + repositoryName,
			"--set-string",
			"labels.commitHash=\"" + commitHash + "\"",
		},
		" ",
	)

	actualCommand := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		repositoryName,
		commitHash,
		false,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmDeployCommand3(t *testing.T) {
	const systemName = "core"
	const helmValuesFile = ".github/deploy/values.yml"
	const applicationName = "demo-api"
	const environment = "prod"
	const workloadType = "job"
	const imageTag = "v420"
	const repositoryName = "core-not-monorepo"
	const commitHash = "abcdef"

	commandOutput := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		repositoryName,
		commitHash,
		false,
		&command.RunOptions{DryRun: true},
	)

	if !command.IsError(commandOutput) {
		t.Errorf("Expected error, got %s", commandOutput)
	}
}
