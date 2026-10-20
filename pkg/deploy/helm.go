package deploy

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/3lvia/cli/pkg/command"
)

const (
	chartsRepositoryURL = "oci://ghcr.io/3lvia/charts"
)

func checkHelmInstalledCommand(
	runOptions *command.RunOptions,
) command.Output {
	return command.Run(
		*exec.Command("helm", "version"),
		runOptions,
	)
}

func chartReference(workloadType string, useISSChart bool) (string, error) {
	switch workloadType {
	case "deployment", "statefulset", "job", "cronjob":
	default:
		return "", fmt.Errorf(
			"Workload type must be either deployment, statefulset, cronjob or job, got '%s'.",
			workloadType,
		)
	}

	if useISSChart {
		if workloadType == "deployment" {
			return chartsRepositoryURL + "/iss-" + workloadType, nil
		}

		return "", fmt.Errorf("Workload type '%s' is not supported with ISS chart.", workloadType)
	}

	return chartsRepositoryURL + "/elvia-" + workloadType, nil
}

// helmValueArgs are the values given to both `helm upgrade` and `helm template`, so that the manifest
// we resolve the image digest from is rendered from exactly the same input as the one we deploy.
func helmValueArgs(
	environment string,
	repositoryName string,
	commitHash string,
	imageTag string,
) []string {
	args := []string{
		"--set-string",
		"environment=" + environment,
		"--set-string",
		"labels.repositoryName=" + repositoryName,
		"--set-string",
		"labels.commitHash=\"" + commitHash + "\"",
	}

	if imageTag != "" {
		args = append(args, "--set-string", "image.tag="+imageTag)
	}

	return args
}

func helmDeployCommand(
	applicationName string,
	systemName string,
	helmValuesFile string,
	environment string,
	workloadType string,
	imageTag string,
	imageDigest string,
	repositoryName string,
	commitHash string,
	dryRun bool,
	useISSChart bool,
	runOptions *command.RunOptions,
) command.Output {
	chart, err := chartReference(workloadType, useISSChart)
	if err != nil {
		return command.Error(err)
	}

	cmd := exec.Command(
		"helm",
		"upgrade",
		"--debug",
		"--install",
		"-n",
		systemName,
		"-f",
		helmValuesFile,
		applicationName,
		chart,
	)
	cmd.Args = append(cmd.Args, helmValueArgs(environment, repositoryName, commitHash, imageTag)...)

	if imageDigest != "" {
		cmd.Args = append(cmd.Args, "--set-string", "image.digest="+imageDigest)
	}

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		cmd.Args = append(cmd.Args, "--set-string", "labels.deployedBy=github-actions")
	}

	if dryRun {
		cmd.Args = append(cmd.Args, "--dry-run")
	}

	return command.Run(*cmd, runOptions)
}

// helmTemplateCommand renders the chart locally with the same input as helmDeployCommand.
// Run it with RunOptions.Quiet: the rendered manifests can contain secret values and must not be logged.
func helmTemplateCommand(
	applicationName string,
	systemName string,
	helmValuesFile string,
	environment string,
	workloadType string,
	imageTag string,
	repositoryName string,
	commitHash string,
	useISSChart bool,
	runOptions *command.RunOptions,
) command.Output {
	chart, err := chartReference(workloadType, useISSChart)
	if err != nil {
		return command.Error(err)
	}

	cmd := exec.Command(
		"helm",
		"template",
		"-n",
		systemName,
		"-f",
		helmValuesFile,
		applicationName,
		chart,
	)
	cmd.Args = append(cmd.Args, helmValueArgs(environment, repositoryName, commitHash, imageTag)...)

	return command.Run(*cmd, runOptions)
}
