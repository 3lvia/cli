package deploy

import (
	"strings"
	"testing"

	"github.com/3lvia/cli/pkg/command"
)

func TestCheckHelmInstalledCommand(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
		"",
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmRepoAddCommandWithUrl(t *testing.T) {
	t.Parallel()

	expectedCommandString := strings.Join(
		[]string{
			"helm",
			"repo",
			"add",
			"elvia-charts",
			"https://raw.githubusercontent.com/3lvia/kubernetes-charts/feature/cool-new-charts",
		},
		" ",
	)

	actualCommand := helmRepoAddCommand(
		"https://raw.githubusercontent.com/3lvia/kubernetes-charts/feature/cool-new-charts",
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmRepoUpdateCommand(t *testing.T) {
	t.Parallel()

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
	t.Setenv("GITHUB_ACTIONS", "false") // Reset GITHUB_ACTIONS env var so tests don't fail in GitHub Actions

	const (
		systemName      = "core"
		helmValuesFile  = ".github/deploy/values.yaml"
		applicationName = "demo-api"
		environment     = "dev"
		workloadType    = "deployment"
		imageTag        = "v12"
		imageDigest     = "sha256:1234567890"
		repositoryName  = "core"
		commitHash      = "123456"
	)

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
			"labels.repositoryName=" + repositoryName,
			"--set-string",
			"labels.commitHash=\"" + commitHash + "\"",
			"--set-string",
			"image.tag=" + imageTag,
			"--set-string",
			"image.digest=" + imageDigest,
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
		imageDigest,
		repositoryName,
		commitHash,
		false,
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
	t.Setenv("GITHUB_ACTIONS", "false") // Reset GITHUB_ACTIONS env var so tests don't fail in GitHub Actions

	const (
		systemName      = "core"
		helmValuesFile  = ".github/deploy/values.yaml"
		applicationName = "demo-api"
		environment     = "prod"
		workloadType    = "statefulset"
		imageTag        = "v420"
		imageDigest     = "sha256:abcdef"
		repositoryName  = "core-not-monorepo"
		commitHash      = "abcdef"
	)

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
			"labels.repositoryName=" + repositoryName,
			"--set-string",
			"labels.commitHash=\"" + commitHash + "\"",
			"--set-string",
			"image.tag=" + imageTag,
			"--set-string",
			"image.digest=" + imageDigest,
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
		imageDigest,
		repositoryName,
		commitHash,
		false,
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
	t.Setenv("GITHUB_ACTIONS", "false") // Reset GITHUB_ACTIONS env var so tests don't fail in GitHub Actions

	const (
		systemName      = "core"
		helmValuesFile  = ".github/deploy/values.yaml"
		applicationName = "demo-api"
		environment     = "prod"
		workloadType    = "job"
		imageTag        = "v420"
		imageDigest     = "sha256:abcdef"
		repositoryName  = "core-not-monorepo"
		commitHash      = "abcdef"
	)

	commandOutput := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		imageDigest,
		repositoryName,
		commitHash,
		false,
		false,
		&command.RunOptions{DryRun: true},
	)

	if !command.IsError(commandOutput) {
		t.Errorf("Expected error, got %s", commandOutput)
	}
}

func TestHelmDeployCommand4(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "false") // Reset GITHUB_ACTIONS env var so tests don't fail in GitHub Actions

	const (
		systemName      = "core"
		helmValuesFile  = ".github/deploy/values.yaml"
		applicationName = "demo-api"
		environment     = "prod"
		workloadType    = "deployment"
		imageTag        = ""
		imageDigest     = "sha256:abcdef"
		repositoryName  = "iss-demo-api"
		commitHash      = "abcdef"
	)

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
			"elvia-charts/iss-" + workloadType,
			"--set-string",
			"environment=" + environment,
			"--set-string",
			"labels.repositoryName=" + repositoryName,
			"--set-string",
			"labels.commitHash=\"" + commitHash + "\"",
			"--set-string",
			"image.digest=" + imageDigest,
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
		imageDigest,
		repositoryName,
		commitHash,
		false,
		true,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}

func TestHelmDeployCommand5(t *testing.T) {
	t.Setenv("GITHUB_ACTIONS", "false") // Reset GITHUB_ACTIONS env var so tests don't fail in GitHub Actions

	const (
		systemName      = "core"
		helmValuesFile  = ".github/deploy/values.yaml"
		applicationName = "demo-api"
		environment     = "prod"
		workloadType    = "statefulset"
		imageTag        = "v420"
		imageDigest     = "sha256:abcdef"
		repositoryName  = "iss-demo-api"
		commitHash      = "abcdef"
	)

	commandOutput := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		imageDigest,
		repositoryName,
		commitHash,
		false,
		true,
		&command.RunOptions{DryRun: true},
	)

	if !command.IsError(commandOutput) {
		t.Errorf("Expected error, got %s", commandOutput)
	}
}

func TestHelmDeployCommandWithGitHubActionsEnv(t *testing.T) {
	const (
		systemName      = "core"
		helmValuesFile  = ".github/deploy/values.yaml"
		applicationName = "demo-api"
		environment     = "dev"
		workloadType    = "deployment"
		imageTag        = "v12"
		imageDigest     = "sha256:1234567890"
		repositoryName  = "core"
		commitHash      = "123456"
	)

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
			"labels.repositoryName=" + repositoryName,
			"--set-string",
			"labels.commitHash=\"" + commitHash + "\"",
			"--set-string",
			"image.tag=" + imageTag,
			"--set-string",
			"image.digest=" + imageDigest,
			"--set-string",
			"labels.deployedBy=github-actions",
		},
		" ",
	)

	t.Setenv("GITHUB_ACTIONS", "true")

	actualCommand := helmDeployCommand(
		applicationName,
		systemName,
		helmValuesFile,
		environment,
		workloadType,
		imageTag,
		imageDigest,
		repositoryName,
		commitHash,
		false,
		false,
		&command.RunOptions{DryRun: true},
	)

	command.ExpectedCommandStringEqualsActualCommand(
		t,
		expectedCommandString,
		actualCommand,
	)
}
